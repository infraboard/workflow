package pipeline

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/infraboard/mcube/http/request"
	"github.com/infraboard/mcube/types/ftime"
)

const (
	AUDIT_NOTIFY_MARK_KEY = "AUDIT_NOTIFY_HAS_SEND"
)

func NewFlow(number int64, items []*Step) *Flow {
	return &Flow{
		number: number,
		items:  items,
	}
}

type Flow struct {
	number int64
	items  []*Step
}

// 判断 这个flow有没有中断
func (f *Flow) IsBreak() bool {
	for i := range f.items {
		step := f.items[i]
		// step break 算stage中断退出
		if step.IsBreakNow() {
			return true
		}
	}
	return false
}

// 判断所有的step是不是都执行成功了
func (f *Flow) IsPassed() bool {
	count := 0
	for i := range f.items {
		step := f.items[i]
		if step.IsPassed() {
			count++
		}
	}
	return count == len(f.items)
}

func (f *Flow) IsComplete() bool {
	return f.IsBreak() || f.IsPassed()
}

func NewDefaultStage() *Stage {
	return &Stage{
		Steps: []*Step{},
	}
}

func (s *Stage) StepCount() int {
	return len(s.Steps)
}

func (s *Stage) ShortDesc() string {
	return fmt.Sprintf("%s[%d]", s.Name, s.Id)
}

func (s *Stage) AddStep(item *Step) {
	s.Steps = append(s.Steps, item)
}

func (s *Stage) UpdateStep(item *Step) error {
	step, err := s.GetStepByKey(item.Key)
	if err != nil {
		return err
	}

	*step = *item
	return nil
}

func (s *Stage) GetStepByKey(key string) (*Step, error) {
	for i := range s.Steps {
		step := s.Steps[i]
		if step.Key == key {
			return step, nil
		}
	}

	return nil, fmt.Errorf("step %s not found", key)
}

func (s *Stage) LastStep() *Step {
	if s.StepCount() == 0 {
		return nil
	}

	return s.Steps[s.StepCount()-1]
}

// 因为可能包含并行任务, 下一次执行的任务可能是多个
func (s *Stage) NextStep() (nextSteps []*Step) {
	for i := range s.Steps {
		step := s.Steps[i]

		// 已经调度的Step不计入下一次调度范围
		if step.IsComplete() {
			continue
		}

		nextSteps = append(nextSteps, step)
		// 遇到串行执行的step结束step
		if !step.IsParallel {
			return
		}
	}
	return
}

func (s *Stage) GetFlow(flowNumber int64) *Flow {
	steps := []*Step{}
	for i := range s.Steps {
		step := s.Steps[i]

		if step.FlowNumber() > flowNumber {
			return nil
		}

		if step.FlowNumber() == flowNumber {
			steps = append(steps, step)
		}
	}

	if len(steps) == 0 {
		return nil
	}

	return NewFlow(flowNumber, steps)
}

func (s *Stage) IsRunning() bool {
	for i := range s.Steps {
		step := s.Steps[i]

		if step.IsRunning() {
			return true
		}
	}
	return false
}

// 判断Stage是否执行成功
func (s *Stage) IsBreakNow() bool {
	for i := range s.Steps {
		step := s.Steps[i]
		if step.IgnoreFailed {
			continue
		}
		if step.IsBreakNow() {
			return true
		}
	}

	return false
}

// 最后一个step都没有中断
func (s *Stage) IsPassed() bool {
	step := s.LastStep()
	if step == nil {
		return true
	}

	return step.IsPassed()
}

// LoadStepFromBytes 解析etcd 的step数据
func LoadStepFromBytes(value []byte) (*Step, error) {
	step := NewDefaultStep()

	// 解析Value
	if len(value) > 0 {
		if err := json.Unmarshal(value, step); err != nil {
			return nil, fmt.Errorf("unmarshal step error, vaule(%s) %s", string(value), err)
		}
	}

	// 校验合法性
	if err := step.Validate(); err != nil {
		return nil, err
	}

	return step, nil
}

// NewStepSet todo
func NewStepSet() *StepSet {
	return &StepSet{
		Items: []*Step{},
	}
}

func (s *StepSet) Add(item *Step) {
	s.Items = append(s.Items, item)
}

func NewDefaultStep() *Step {
	return &Step{
		Status: NewDefaultStepStatus(),
	}
}

func (s *Step) Run() {
	s.Status.StartAt = ftime.Now().Timestamp()
	s.Status.Status = STEP_STATUS_RUNNING
}

func (s *Step) FlowNumber() int64 {
	if s.Status == nil {
		return 0
	}

	return s.Status.FlowNumber
}

func (s *Step) setFlowNumber(n int64) {
	s.Status.FlowNumber = n
}

func (s *Step) Failed(format string, a ...interface{}) {
	s.Status.EndAt = ftime.Now().Timestamp()
	s.Status.Status = STEP_STATUS_FAILED
	s.Status.Message = fmt.Sprintf(format, a...)
}

func (s *Step) Cancel(format string, a ...interface{}) {
	s.Status.Status = STEP_STATUS_CANCELING
	s.Status.Message = fmt.Sprintf(format, a...)
}

func (s *Step) Audit(resp AUDIT_RESPONSE, message string) {
	s.Status.AuditAt = ftime.Now().Timestamp()
	s.Status.AuditResponse = resp
	s.Status.AuditMessage = message
}

func (s *Step) HasSendAuditNotify() bool {
	if s.Status.ContextMap == nil {
		return false
	}

	return s.Status.ContextMap[AUDIT_NOTIFY_MARK_KEY] == "true"
}

func (s *Step) MarkSendAuditNotify() {
	if s.Status.ContextMap == nil {
		s.Status.ContextMap = map[string]string{}
	}

	s.Status.ContextMap[AUDIT_NOTIFY_MARK_KEY] = "true"
	s.Status.Status = STEP_STATUS_AUDITING
}

func (s *Step) Success(resp map[string]string) {
	s.Status.EndAt = ftime.Now().Timestamp()
	s.Status.Status = STEP_STATUS_SUCCEEDED
	s.UpdateResponse(resp)
}

func (s *Step) UpdateResponse(resp map[string]string) {
	if s.Status.Response == nil {
		s.Status.Response = map[string]string{}
	}
	for k, v := range resp {
		s.Status.Response[k] = v
	}
}

func (s *Step) MakeObjectKey() string {
	return StepObjectKey(s.Key)
}

func (s *Step) Validate() error {
	return validate.Struct(s)
}

func (s *Step) SetScheduleNode(nodeName string) {
	s.Status.ScheduledNode = nodeName
}

func (s *Step) ScheduledNodeName() string {
	if s.Status != nil {
		return s.Status.ScheduledNode
	}
	return ""
}

func (s *Step) IsScheduled() bool {
	return s.ScheduledNodeName() != ""
}

func (s *Step) IsComplete() bool {
	if s.Status == nil {
		return false
	}

	return s.Status.Status.IsIn(
		STEP_STATUS_SUCCEEDED,
		STEP_STATUS_FAILED,
		STEP_STATUS_CANCELED,
		STEP_STATUS_SKIP,
		STEP_STATUS_REFUSE,
	)
}

func (s *Step) IsRunning() bool {
	if s.Status == nil {
		return false
	}

	return s.Status.Status.IsIn(
		STEP_STATUS_RUNNING,
		STEP_STATUS_CANCELING,
		STEP_STATUS_AUDITING,
	)
}

func (s *Step) IsBreakNow() bool {
	if s.Status == nil {
		return false
	}

	// 忽略失败的已计算为通过
	if s.IgnoreFailed {
		return false
	}

	return s.Status.Status.IsIn(
		STEP_STATUS_FAILED,
		STEP_STATUS_CANCELED,
		STEP_STATUS_CANCELING,
		STEP_STATUS_REFUSE,
	)
}

func (s *Step) IsPassed() bool {
	if s.Status == nil {
		return false
	}

	// 忽略失败的已计算为通过
	if s.IgnoreFailed {
		return true
	}

	return s.Status.Status.IsIn(
		STEP_STATUS_SUCCEEDED,
		STEP_STATUS_SKIP,
	)
}

func (s *Step) BuildKey(namespace, pipelineId string, stage int32) {
	s.Key = fmt.Sprintf("%s.%s.%d.%d", namespace, pipelineId, stage, s.Id)
}

func (s *Step) AuditPass() bool {
	return s.Status.AuditResponse.Equal(AUDIT_RESPONSE_ALLOW)
}

func (s *Step) GetPipelineStepNumber() int32 {
	n, _ := strconv.ParseInt(s.getKeyIndex(3), 10, 32)
	return int32(n)
}

func (s *Step) GetPipelineStageNumber() int32 {
	n, _ := strconv.ParseInt(s.getKeyIndex(2), 10, 32)
	return int32(n)
}

func (s *Step) GetPipelineID() string {
	return s.getKeyIndex(1)
}

func (s *Step) GetNamespace() string {
	return s.getKeyIndex(0)
}

func (s *Step) getKeyIndex(index int) string {
	kl := strings.Split(s.Key, ".")
	if index+1 > len(kl) {
		return ""
	}

	if len(kl) != 4 {
		return ""
	}

	return kl[index]
}

func NewDefaultStepStatus() *StepStatus {
	return &StepStatus{
		Response: map[string]string{},
	}
}

// NewQueryStepRequest 查询book列表
func NewQueryStepRequest() *QueryStepRequest {
	return &QueryStepRequest{
		Page: &request.NewDefaultPageRequest().PageRequest,
	}
}

// NewDescribeStepRequestWithKey 查询book列表
func NewDescribeStepRequestWithKey(key string) *DescribeStepRequest {
	return &DescribeStepRequest{
		Key: key,
	}
}

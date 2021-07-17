package pipeline

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/infraboard/mcube/http/request"
	"github.com/infraboard/mcube/types/ftime"
)

func NewDefaultStage() *Stage {
	return &Stage{
		Steps: []*Step{},
	}
}

func (s *Stage) StepCount() int {
	return len(s.Steps)
}

func (s *Stage) AddStep(item *Step) {
	s.Steps = append(s.Steps, item)
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
		if step.IsScheduled() || step.IsComplete() {
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

// 判断Stage最后一个任务是否完成
func (s *Stage) IsComplete() bool {
	step := s.LastStep()
	if step == nil {
		return true
	}

	return step.IsComplete()
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

func (s *Step) Failed(format string, a ...interface{}) {
	s.Status.EndAt = ftime.Now().Timestamp()
	s.Status.Status = STEP_STATUS_FAILED
	s.Status.Message = fmt.Sprintf(format, a...)
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

func (s *Step) IsAudit() bool {
	return s.Status.AuditAt != 0
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
	if s.Status != nil {
		return s.Status.Status.IsIn(
			STEP_STATUS_SUCCEEDED,
			STEP_STATUS_FAILED,
			STEP_STATUS_CANCELED,
			STEP_STATUS_SKIP,
			STEP_STATUS_REFUSE,
		)
	}

	return false
}

func (s *Step) BuildKey(namespace, pipelineId string, stage int32) {
	s.Key = fmt.Sprintf("%s.%s.%d.%d", namespace, pipelineId, stage, s.Id)
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

package pipeline

import (
	"encoding/json"
	"fmt"
	"regexp"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/infraboard/keyauth/app/token"
	"github.com/infraboard/mcube/http/request"
	"github.com/infraboard/mcube/logger/zap"
	"github.com/rs/xid"
)

const (
	// 一个pipeline最多可以有少个step
	PIPELINE_MAX_STEPS = 100
)

var (
	VALUE_TYPE_ID_MAP = map[string]PARAM_VALUE_TYPE{
		"$p$": PARAM_VALUE_TYPE_PASSWORD,
		"$c$": PARAM_VALUE_TYPE_CRYPTO,
		"$a$": PARAM_VALUE_TYPE_APP_VAR,
		"$s$": PARAM_VALUE_TYPE_SECRET_REF,
	}
)

// use a single instance of Validate, it caches struct info
var (
	validate = validator.New()
)

func NewCreatePipelineRequest() *CreatePipelineRequest {
	return &CreatePipelineRequest{}
}

func (req *CreatePipelineRequest) UpdateOwner(tk *token.Token) {
	req.Domain = tk.Domain
	req.Namespace = tk.Namespace
	req.CreateBy = tk.Account
}

// NewQueryPipelineRequest 查询book列表
func NewQueryPipelineRequest() *QueryPipelineRequest {
	return &QueryPipelineRequest{
		Page: request.NewDefaultPageRequest(),
	}
}

func (req *CreatePipelineRequest) Validate() error {
	if len(req.Stages) == 0 {
		return fmt.Errorf("no stages")
	}
	return validate.Struct(req)
}

func (req *CreatePipelineRequest) EnsureStep() {
	for m := range req.Stages {
		s := req.Stages[m]
		s.Id = int32(m) + 1
		for n := range s.Steps {
			t := s.Steps[n]
			t.Id = int32(n) + 1
			t.Status = NewDefaultStepStatus()
		}
	}
}

func LoadPipelineFromBytes(payload []byte) (*Pipeline, error) {
	ins := NewDefaultPipeline()

	// 解析Value
	if err := json.Unmarshal(payload, ins); err != nil {
		return nil, fmt.Errorf("unmarshal step error, vaule(%s) %s", string(payload), err)
	}

	// 校验合法性
	if err := ins.Validate(); err != nil {
		return nil, err
	}

	return ins, nil
}

func NewDefaultPipeline() *Pipeline {
	return &Pipeline{
		Status: &PipelineStatus{},
	}
}

func NewDefaultPipelineStatus() *PipelineStatus {
	return &PipelineStatus{}
}

func NewPipeline(req *CreatePipelineRequest) (*Pipeline, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	// 补充step id
	req.EnsureStep()

	p := &Pipeline{
		Id:          xid.New().String(),
		CreateAt:    time.Now().UnixMilli(),
		TemplateId:  req.TemplateId,
		CreateBy:    req.CreateBy,
		Domain:      req.Domain,
		Namespace:   req.Namespace,
		Name:        req.Name,
		With:        req.With,
		Mount:       req.Mount,
		Tags:        req.Tags,
		Description: req.Description,
		On:          req.On,
		Stages:      req.Stages,
		HookEvent:   req.HookEvent,
		Status:      &PipelineStatus{},
	}
	return p, nil
}

// MakeObjectKey 构建etcd对应的key
// 例如: inforboard/workflow/service/node/node-01
func (p *Pipeline) MakeObjectKey() string {
	return fmt.Sprintf("%s/%s/%s", EtcdPipelinePrefix(), p.Namespace, p.Id)
}

func (p *Pipeline) UpdateOwner(tk *token.Token) {
	p.CreateBy = tk.Account
	p.Domain = tk.Domain
	p.Namespace = tk.Namespace
}

func (p *Pipeline) Validate() error {
	return validate.Struct(p)
}

func (p *Pipeline) GetStep(stageNumber int32, stepKey string) (*Step, error) {
	stage, err := p.GetStageByNumber(stageNumber)
	if err != nil {
		return nil, fmt.Errorf("get stage error, %s", err)
	}

	return stage.GetStepByKey(stepKey)
}

func (p *Pipeline) UpdateStep(s *Step) error {
	ns, id := s.GetNamespace(), s.GetPipelineId()
	if ns != p.Namespace || id != p.Id {
		return fmt.Errorf("this step not match this pipeline, id or namespace is not correct")
	}

	stage, err := p.GetStageByNumber(s.GetPipelineStageNumber())
	if err != nil {
		return fmt.Errorf("get stage error, %s", err)
	}

	return stage.UpdateStep(s)
}

func (p *Pipeline) AddStage(item *Stage) {
	p.Stages = append(p.Stages, item)
}

func (p *Pipeline) GetStageByNumber(number int32) (*Stage, error) {
	if int(number) > len(p.Stages) || number <= 0 {
		return nil, fmt.Errorf("number range 1 ~ %d", len(p.Stages))
	}

	return p.Stages[number-1], nil
}

func (p *Pipeline) LastStage() *Stage {
	total := len(p.Stages)
	if total == 0 {
		return nil
	}

	return p.Stages[total-1]
}

func (p *Pipeline) IsComplete() bool {
	return p.Status.Status.Equal(PIPELINE_STATUS_COMPLETE)
}

func (s *Pipeline) IsScheduled() bool {
	if s.Status == nil {
		return false
	}

	return s.Status.SchedulerNode != ""
}

func (p *Pipeline) IsRunning() bool {
	if p.Status == nil {
		return false
	}
	return p.Status.Status.Equal(PIPELINE_STATUS_EXECUTING)
}

func (p *Pipeline) NextFlowNumber() int64 {
	if p.Status == nil {
		return 1
	}

	return p.Status.GetCurrentFlow() + 1
}

func (p *Pipeline) CurrentFlowNumber() int64 {
	if p.Status == nil {
		return 0
	}

	return p.Status.GetCurrentFlow()
}

func (p *Pipeline) incFlow() {
	p.Status.CurrentFlow++
}

func (p *Pipeline) Run() {
	if p.Status == nil {
		p.Status = NewDefaultPipelineStatus()
	}
	p.Status.Status = PIPELINE_STATUS_EXECUTING
	p.Status.StartAt = time.Now().UnixMilli()
}

func (p *Pipeline) Complete() {
	if p.Status == nil {
		p.Status = NewDefaultPipelineStatus()
	}
	p.Status.Status = PIPELINE_STATUS_COMPLETE
	p.Status.EndAt = time.Now().UnixMilli()
}

func (p *Pipeline) GetCurrentFlow() *Flow {
	flow := p.CurrentFlowNumber()

	// pipeline还没有运行
	if flow == 0 {
		return nil
	}

	for i := range p.Stages {
		stage := p.Stages[i]
		if f := stage.GetFlow(flow); f != nil {
			return f
		}
	}

	return nil
}

func (p *Pipeline) GetNextFlow() *Flow {
	for i := range p.Stages {
		stage := p.Stages[i]

		// stage 通过了继续搜索下一个stage
		if stage.IsPassed() {
			continue
		}

		// 如果stage中断后 就没有下一步了
		if stage.IsBreakNow() {
			return nil
		}

		steps := stage.NextStep()
		for i := range steps {
			step := steps[i]
			step.PipelineId = p.Id
			step.Namespace = p.Namespace
			step.CreateAt = time.Now().UnixMilli()
			steps[i].BuildKey(p.Namespace, p.Id, stage.Id)
			step.setFlowNumber(p.NextFlowNumber())
		}

		if len(steps) > 0 {
			p.incFlow()
		}
		return NewFlow(p.CurrentFlowNumber(), steps)
	}
	return nil
}

// 只有上一个flow执行完成后, 才会有下个fow
// 注意: 多个并行的任务是不能跨stage同时执行的
//      也就是说stage一定是串行的
func (p *Pipeline) NextStep() (steps []*Step, isComplete bool) {
	// 判断当前Flow是否运行完成
	if f := p.GetCurrentFlow(); f != nil {
		// 如果有flow中断, pipeline 提前结束
		if f.IsBreak() {
			isComplete = true
			return
		}

		// 如果flow没有pass 说明还是在运行中, 不需要调度下以组step
		if !f.IsPassed() {
			return
		}
	}

	f := p.GetNextFlow()

	// 判断是不是最后一个Flow了
	if f == nil {
		isComplete = true
		return
	}

	// 如果不是则获取flow中的step
	steps = f.items
	return
}

func (p *Pipeline) ShortDescribe() string {
	return fmt.Sprintf("%s[%s]", p.Name, p.Id)
}

func (p *Pipeline) ScheduledNodeName() string {
	return p.Status.SchedulerNode
}

func (p *Pipeline) SetScheduleNode(nodeName string) {
	p.Status.SchedulerNode = nodeName
}

func (p *Pipeline) EtcdObjectKey() string {
	return fmt.Sprintf("%s/%s/%s", EtcdPipelinePrefix(), p.Namespace, p.Id)
}

func (p *Pipeline) MatchScheduler(schedulerName string) bool {
	if p.Status == nil {
		return false
	}
	return p.Status.SchedulerNode == schedulerName
}

func (p *Pipeline) StepPrefix() string {
	return fmt.Sprintf("%s.%s", p.Namespace, p.Id)
}

// NewPipelineSet todo
func NewPipelineSet() *PipelineSet {
	return &PipelineSet{
		Items: []*Pipeline{},
	}
}

func (s *PipelineSet) Add(item *Pipeline) {
	s.Items = append(s.Items, item)
}

// NewQueryPipelineRequest 查询book列表
func NewDescribePipelineRequestWithID(id string) *DescribePipelineRequest {
	return &DescribePipelineRequest{
		Id: id,
	}
}

// NewDeletePipelineRequestWithID 查询book列表
func NewDeletePipelineRequestWithID(id string) *DeletePipelineRequest {
	return &DeletePipelineRequest{
		Id: id,
	}
}

func NewQueryPipelineOptions() *QueryPipelineOptions {
	return &QueryPipelineOptions{}
}

// QueryPipelineTaskOptions ListPipeline 查询条件
type QueryPipelineOptions struct {
	Node string
}

func NewDeleteStepRequestWithKey(key string) *DeleteStepRequest {
	return &DeleteStepRequest{
		Key: key,
	}
}

func NewAuditStepRequest() *AuditStepRequest {
	return &AuditStepRequest{}
}

func NewAuditStepRequestWithKey(key string) *AuditStepRequest {
	return &AuditStepRequest{
		Key: key,
	}
}

func NewCancelStepRequestWithKey(key string) *CancelStepRequest {
	return &CancelStepRequest{
		Key: key,
	}
}

func NewWatchPipelineRequestByID(namespace, id string) *CreateWatchPipelineRequest {
	return &CreateWatchPipelineRequest{
		Namespace: namespace,
		Id:        id,
		Mod:       PIPELINE_WATCH_MOD_BY_ID,
	}
}

func NewWatchPipelineRequestByNamespace(namespace string) *CreateWatchPipelineRequest {
	return &CreateWatchPipelineRequest{
		Namespace: namespace,
		Mod:       PIPELINE_WATCH_MOD_BY_NAMESPACE,
	}
}

func (req *CreateWatchPipelineRequest) Validate() error {
	switch req.Mod {
	case PIPELINE_WATCH_MOD_BY_ID:
		if req.Id != "" && req.Namespace == "" {
			return fmt.Errorf("when watch pipeline namespace id required")
		}
	case PIPELINE_WATCH_MOD_BY_NAMESPACE:
		if req.Namespace == "" {
			return fmt.Errorf("namespace required")
		}
	default:
		return fmt.Errorf("unknown watch mod %s", req.Mod)
	}

	return nil
}

func (t *Trigger) IsMatch(branche, event string) bool {
	return t.matchBranche(branche) && t.matchEvent(event)
}

func (t *Trigger) matchBranche(b string) bool {
	for i := range t.Branches {
		matched, err := regexp.MatchString(b, t.Branches[i])
		if err != nil {
			zap.L().Errorf("match branche string error, %s", err)
		}
		if matched {
			return true
		}
	}
	return false
}

func (t *Trigger) matchEvent(e string) bool {
	for i := range t.Events {
		matched, err := regexp.MatchString(e, t.Events[i])
		if err != nil {
			zap.L().Errorf("match event string error, %s", err)
		}
		if matched {
			return true
		}
	}

	return false
}

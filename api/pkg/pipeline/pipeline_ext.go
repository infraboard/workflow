package pipeline

import (
	"encoding/json"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/infraboard/mcube/http/request"
	"github.com/infraboard/mcube/types/ftime"
	"github.com/rs/xid"
)

// use a single instance of Validate, it caches struct info
var (
	validate = validator.New()
)

func NewCreatePipelineRequest() *CreatePipelineRequest {
	return &CreatePipelineRequest{}
}

func NewDefaultStepStatus() *StepStatus {
	return &StepStatus{}
}

// NewQueryPipelineRequest 查询book列表
func NewQueryPipelineRequest() *QueryPipelineRequest {
	return &QueryPipelineRequest{
		Page: &request.NewDefaultPageRequest().PageRequest,
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

func NewPipeline(req *CreatePipelineRequest) (*Pipeline, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	// 补充step id
	req.EnsureStep()

	p := &Pipeline{
		Id:          xid.New().String(),
		CreateAt:    ftime.Now().Timestamp(),
		Name:        req.Name,
		Tags:        req.Tags,
		Description: req.Description,
		On:          req.On,
		Stages:      req.Stages,
		Status:      &PipelineStatus{},
	}
	return p, nil
}

func (p *Pipeline) Validate() error {
	return validate.Struct(p)
}

func (t *Pipeline) NextStep() (steps []*Step) {
	for i := range t.Stages {
		stage := t.Stages[i]
		steps = stage.NextStep()
		for i := range steps {
			steps[i].BuildKey(t.Id, stage.Id)
		}
	}

	return
}

func (t *Pipeline) ShortDescribe() string {
	return fmt.Sprintf("%s[%s]", t.Name, t.Id)
}

func (t *Pipeline) SchedulerNodeName() string {
	return t.Status.SchedulerNode
}

func (t *Pipeline) SetScheduleNode(nodeName string) {
	t.Status.SchedulerNode = nodeName
}

func (p *Pipeline) EtcdObjectKey() string {
	return fmt.Sprintf("%s/%s/%s", EtcdPipelinePrefix(), p.Namespace, p.Id)
}

func (s *PipelineStatus) IsScheduled() bool {
	return s.SchedulerNode != ""
}

func (s *PipelineStatus) MatchScheduler(schedulerName string) bool {
	return s.SchedulerNode == schedulerName
}

func (s *PipelineStatus) IsComplete() bool {
	return s.Status.Equal(PIPELINE_STATUS_COMPLETE)
}

func (s *PipelineStatus) IsRunning() bool {
	return s.Status.Equal(PIPELINE_STATUS_EXECUTING)
}

func (s *PipelineStatus) Run() {
	s.Status = PIPELINE_STATUS_EXECUTING
	if s.StartAt != 0 {
		s.StartAt = ftime.Now().Timestamp()
	}
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

func (s *Stage) StepCount() int {
	return len(s.Steps)
}

// 因为可能包含并行任务, 下一次执行的任务可能是多个
func (s *Stage) NextStep() (nextSteps []*Step) {
	for i := range s.Steps {
		step := s.Steps[i]

		// 已经调度的Step不计入下一次调度范围
		if step.IsScheduled() {
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

func NewDefaultStep() *Step {
	return &Step{}
}

func (s *Step) Validate() error {
	return validate.Struct(s)
}

func (s *Step) EtcdObjectKey(prefix string) string {
	return fmt.Sprintf("%s/%s", prefix, s.Key)
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

func (s *Step) BuildKey(pipelineId string, stage int32) {
	s.Key = fmt.Sprintf("%s.%d.%d", pipelineId, stage, s.Id)
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

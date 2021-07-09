package pipeline

import (
	"encoding/json"
	"fmt"

	"github.com/infraboard/mcube/http/request"
)

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
	return &Step{}
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

func (s *Step) BuildKey(namespace, pipelineId string, stage int32) {
	s.Key = fmt.Sprintf("%s.%s.%d.%d", namespace, pipelineId, stage, s.Id)
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

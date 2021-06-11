package task

import (
	"encoding/json"
	"fmt"

	"github.com/infraboard/workflow/api/pkg/pipeline"
)

// LoadPipelineTaskFromBytes 解析etcd 的pipeline Task数据
func LoadPipelineTaskFromBytes(value []byte) (*PipelineTask, error) {
	p := NewDefaultPipelineTask()

	// 解析Value
	if len(value) > 0 {
		if err := json.Unmarshal(value, p); err != nil {
			return nil, fmt.Errorf("unmarshal pipeline error, vaule(%s) %s", string(value), err)
		}
	}

	// 校验合法性
	if err := p.Validate(); err != nil {
		return nil, err
	}

	return p, nil
}

func NewDefaultPipelineTask() *PipelineTask {
	return &PipelineTask{}
}

func (t *PipelineTask) Validate() error {
	return nil
}

func (t *PipelineTask) SchedulerNodeName() string {
	return t.SchedulerNode
}

func (t *PipelineTask) AddScheduleNode(nodeName string) {
	t.SchedulerNode = nodeName
}

func (t *PipelineTask) EtcdObjectKey(prefix string) string {
	return fmt.Sprintf("%s/%s", prefix, t.Id)
}

func (t *PipelineTask) NextStep() (steps []*pipeline.Step) {
	for i := range t.Pipeline.Stages {
		stage := t.Pipeline.Stages[i]
		if !stage.HasNextStep() {
			continue
		}

		steps = stage.NextStep()
		for i := range steps {
			steps[i].BuildKey(t.Id, stage.Id)
		}
		return
	}

	return
}

// NewPipelineTaskSet todo
func NewPipelineTaskSet() *PipelineTaskSet {
	return &PipelineTaskSet{
		Items: []*PipelineTask{},
	}
}

func (s *PipelineTaskSet) Add(item *PipelineTask) {
	s.Items = append(s.Items, item)
}

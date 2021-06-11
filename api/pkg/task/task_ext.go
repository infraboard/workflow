package task

import (
	"encoding/json"
	"fmt"
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

// NewPipelineTaskSet todo
func NewPipelineTaskSet() *PipelineTaskSet {
	return &PipelineTaskSet{
		Items: []*PipelineTask{},
	}
}

func (s *PipelineTaskSet) Add(item *PipelineTask) {
	s.Items = append(s.Items, item)
}

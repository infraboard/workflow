package pipeline

import (
	"encoding/json"
	"fmt"

	"github.com/go-playground/validator/v10"
)

// use a single instance of Validate, it caches struct info
var (
	validate = validator.New()
)

func (p *Pipeline) Validate() error {
	return validate.Struct(p)
}

func NewDefaultPipeline() *Pipeline {
	return &Pipeline{}
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

func (p *Pipeline) EtcdObjectKey(prefix string) string {
	return fmt.Sprintf("%s/%s", prefix, p.Id)
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
	return fmt.Sprintf("%s/%d", prefix, s.Id)
}

func (s *Step) AddScheduleNode(nodeName string) {
	s.Status.ScheduledNode = nodeName
}

func (s *Step) ScheduledNodeName() string {
	return ""
}

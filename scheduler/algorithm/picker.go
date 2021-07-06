package algorithm

import (
	"github.com/infraboard/workflow/api/pkg/node"
	"github.com/infraboard/workflow/api/pkg/pipeline"
)

// Picker 挑选一个合适的node 运行Step
type StepPicker interface {
	Pick(*pipeline.Step) (*node.Node, error)
}

type PipelinePicker interface {
	Pick(*pipeline.Pipeline) (*node.Node, error)
}

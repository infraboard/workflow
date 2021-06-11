package algorithm

import (
	"github.com/infraboard/workflow/api/pkg/node"
	"github.com/infraboard/workflow/api/pkg/pipeline"
	"github.com/infraboard/workflow/api/pkg/task"
)

// Picker 挑选一个合适的node 运行Step
type StepPicker interface {
	Pick(*pipeline.Step) (*node.Node, error)
}

type TaskPicker interface {
	Pick(*task.PipelineTask) (*node.Node, error)
}

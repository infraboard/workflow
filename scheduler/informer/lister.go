package informer

import (
	"context"

	"github.com/infraboard/workflow/api/pkg/node"
	"github.com/infraboard/workflow/api/pkg/pipeline"
	"github.com/infraboard/workflow/api/pkg/task"
)

// Lister 所有的Lister
type Lister interface {
	PipelineTaskLister
	StepUpdater
	NodeLister
}

// QueryPipelineTaskOptions ListPipeline 查询条件
type QueryPipelineTaskOptions struct {
	Node string
}

// PipelineTaskLister 获取所有执行节点
type PipelineTaskLister interface {
	ListPipelineTask(ctx context.Context, opts *QueryPipelineTaskOptions) (*task.PipelineTaskSet, error)
}

// StepUpdater todo
type StepUpdater interface {
	UpdateStep(*pipeline.Step) error
}

// NodeLister 获取所有执行节点
type NodeLister interface {
	// List lists all Node
	ListNode(ctx context.Context) ([]*node.Node, error)
}

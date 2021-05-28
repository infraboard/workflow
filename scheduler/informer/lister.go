package informer

import (
	"context"

	"github.com/infraboard/workflow/api/pkg/node"
	"github.com/infraboard/workflow/api/pkg/pipeline"
)

// Lister 所有的Lister
type Lister interface {
	PipelineLister
	StepUpdater
	NodeLister
}

// QueryPipelineOptions ListPipeline 查询条件
type QueryPipelineOptions struct {
	Node string
}

// PipelineLister 获取所有执行节点
type PipelineLister interface {
	ListPipeline(ctx context.Context, opts *QueryPipelineOptions) (*pipeline.PipelineSet, error)
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

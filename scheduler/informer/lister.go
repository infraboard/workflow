package informer

import (
	"context"

	"github.com/infraboard/workflow/api/pkg/node"
	"github.com/infraboard/workflow/api/pkg/pipeline"
)

// Lister 所有的Lister
type Lister interface {
	PipelineLister
	PipelineUpdater
	StepUpdater
	NodeLister
}

func NewQueryPipelineOptions() *QueryPipelineOptions {
	return &QueryPipelineOptions{}
}

// QueryPipelineTaskOptions ListPipeline 查询条件
type QueryPipelineOptions struct {
	Node string
}

// PipelineTaskLister 获取所有执行节点
type PipelineLister interface {
	ListPipeline(ctx context.Context, opts *QueryPipelineOptions) (*pipeline.PipelineSet, error)
}

// StepUpdater todo
type StepUpdater interface {
	UpdateStep(*pipeline.Step) error
}

// StepUpdater todo
type PipelineUpdater interface {
	UpdatePipeline(*pipeline.Pipeline) error
}

// NodeLister 获取所有执行节点
type NodeLister interface {
	// List lists all Node
	ListNode(ctx context.Context) ([]*node.Node, error)
}

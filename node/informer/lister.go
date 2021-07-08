package informer

import (
	"context"

	"github.com/infraboard/workflow/api/pkg/pipeline"
)

// Lister 所有的Lister
type Lister interface {
	StepUpdater
	StepLister
}

// StepUpdater todo
type StepUpdater interface {
	UpdateStep(*pipeline.Step) error
}

// StepLister 获取所有执行节点
type StepLister interface {
	ListStep(ctx context.Context) ([]*pipeline.Step, error)
}

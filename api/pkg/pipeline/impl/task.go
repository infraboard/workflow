package impl

import (
	"context"

	"github.com/infraboard/workflow/api/pkg/pipeline"
)

func (i *impl) CreatePipeline(context.Context, *pipeline.CreatePipelineRequest) (
	*pipeline.Pipeline, error) {

	// 写入key
	// if _, err := i.client.Put(context.Background(), key, value); err != nil {
	// 	return nil, fmt.Errorf("registe service '%s' with ttl to etcd3 failed: %s", key, err.Error())
	// }
	return nil, nil
}

func (i *impl) QueryPipeline(context.Context, *pipeline.QueryPipelineRequest) (
	*pipeline.PipelineSet, error) {
	return nil, nil
}

func (i *impl) CreateAction(context.Context, *pipeline.CreateActionRequest) (
	*pipeline.Action, error) {
	return nil, nil
}

func (i *impl) QueryAction(context.Context, *pipeline.QueryActionRequest) (
	*pipeline.ActionSet, error) {
	return nil, nil
}

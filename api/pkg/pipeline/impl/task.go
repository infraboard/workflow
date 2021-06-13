package impl

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/infraboard/workflow/api/pkg/pipeline"
)

func (i *impl) CreatePipeline(ctx context.Context, req *pipeline.CreatePipelineRequest) (
	*pipeline.Pipeline, error) {

	p, err := pipeline.NewPipeline(req)
	if err != nil {
		return nil, err
	}

	value, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}

	fmt.Println(value)

	// if _, err := i.client.Put(context.Background(), p.EtcdObjectKey(), string(value)); err != nil {
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

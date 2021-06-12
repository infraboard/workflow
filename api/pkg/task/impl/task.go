package impl

import (
	"context"

	"github.com/infraboard/workflow/api/pkg/task"
)

func (i *impl) RunPipeline(ctx context.Context, req *task.RunPipelineRequest) (
	*task.PipelineTask, error) {

	// 写入key
	// if _, err := i.client.Put(context.Background(), key, value); err != nil {
	// 	return nil, fmt.Errorf("registe service '%s' with ttl to etcd3 failed: %s", key, err.Error())
	// }
	return nil, nil
}

func (i *impl) QueryPipelineTask(context.Context, *task.QueryPipelineTaskRequest) (
	*task.PipelineTaskSet, error) {
	return nil, nil
}

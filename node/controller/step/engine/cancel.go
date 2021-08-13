package engine

import (
	"context"

	"github.com/infraboard/mcube/grpc/gcontext"

	"github.com/infraboard/workflow/api/pkg/pipeline"
	"github.com/infraboard/workflow/node/controller/step/runner"
)

func (e *Engine) CancelStep(s *pipeline.Step) {
	if !e.init {
		s.Failed("engine not init")
		return
	}

	e.log.Debugf("start cancel step: %s", s.Key)
	// 构造运行请求
	req := runner.NewCancelRequest(s)

	// 1.查询step对应的action定义
	descA := pipeline.NewDescribeActionRequest(s.GetNamespace(), s.ActionName(), s.ActionVersion())
	ctx := gcontext.NewGrpcOutCtx()
	action, err := e.wc.Pipeline().DescribeAction(ctx.Context(), descA)
	if err != nil {
		s.Failed("describe step action error, %s", err)
		return
	}

	// 3.根据action定义的runner_type, 调用具体的runner
	switch action.RunnerType {
	case pipeline.RUNNER_TYPE_DOCKER:
		go e.docker.Cancel(context.Background(), req)
	case pipeline.RUNNER_TYPE_K8s:
		go e.k8s.Cancel(context.Background(), req)
	case pipeline.RUNNER_TYPE_LOCAL:
		go e.local.Cancel(context.Background(), req)
	default:
		s.Failed("unknown runner type: %s", action.RunnerType)
		return
	}
}

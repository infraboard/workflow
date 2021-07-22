package engine

import (
	"context"

	"github.com/infraboard/mcube/grpc/gcontext"

	"github.com/infraboard/workflow/api/pkg/pipeline"
	"github.com/infraboard/workflow/node/controller/step/runner"
)

// Run 运行Step
// step的参数加载优先级:
//   1. step 本身传人的
//   2. pipeline 运行中产生的
//   3. pipeline 全局传人
//   4. action 默认默认值
func (e *Engine) Run(s *pipeline.Step) {
	if !e.init {
		s.Failed("engine not init")
		return
	}

	e.log.Debugf("start run step: %s", s.Key)
	// 构造运行请求
	req := runner.NewRunRequest(s, e.updateStep)

	// 1.查询step对应的action定义
	descA := pipeline.NewDescribeActionRequestWithName(s.Action)
	descA.Namespace = s.GetNamespace()
	ctx := gcontext.NewGrpcOutCtx()
	action, err := e.wc.Pipeline().DescribeAction(ctx.Context(), descA)
	if err != nil {
		s.Failed("describe step action error, %s", err)
		return
	}

	// 2.查询Pipeline, 获取全局参数
	descP := pipeline.NewDescribePipelineRequestWithID(s.GetPipelineID())
	descP.Namespace = s.GetNamespace()
	pl, err := e.wc.Pipeline().DescribePipeline(ctx.Context(), descP)
	if err != nil {
		s.Failed("describe step pipeline error, %s", err)
		return
	}

	// 加载运行时参数
	req.LoadRunParams(action.DefaultRunParam())
	req.LoadRunParams(pl.With)
	req.LoadMount(pl.Mount)

	// 校验run参数合法性
	if err := action.ValidateRunParam(req.RunParams); err != nil {
		s.Failed(err.Error())
		return
	}

	// 加载Runner运行需要的参数
	req.LoadRunnerParams(action.DefaultRunnerParam())

	// 校验runner参数合法性
	if err := action.ValidateRunnerParam(req.RunnerParams); err != nil {
		s.Failed(err.Error())
		return
	}

	// 3.根据action定义的runner_type, 调用具体的runner
	switch action.RunnerType {
	case pipeline.RUNNER_TYPE_DOCKER:
		go e.docker.Run(context.Background(), req)
	case pipeline.RUNNER_TYPE_K8s:
		go e.k8s.Run(context.Background(), req)
	case pipeline.RUNNER_TYPE_LOCAL:
		go e.local.Run(context.Background(), req)
	default:
		s.Failed("unknown runner type: %s", action.RunnerType)
		return
	}
}

// 如果step执行完成
func (e *Engine) updateStep(s *pipeline.Step) {
	e.log.Debugf("receive step %s update, status %s", s.Key, s.Status)

	if err := e.recorder.Update(s); err != nil {
		e.log.Errorf("update step status error, %s", err)
	}
}

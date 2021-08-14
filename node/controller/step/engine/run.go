package engine

import (
	"context"

	"github.com/infraboard/mcube/grpc/gcontext"

	"github.com/infraboard/workflow/api/pkg/pipeline"
	"github.com/infraboard/workflow/node/controller/step/runner"
)

func (e *Engine) Run(s *pipeline.Step) {
	req := runner.NewRunRequest(s)
	resp := runner.NewRunReponse(e.updateStep)

	e.run(req, resp)

	if resp.HasError() {
		s.Failed(resp.ErrorMessage())
	} else {
		s.Success("")
	}

	e.updateStep(s)
}

// Run 运行Step
// step的参数加载优先级:
//   1. step 本身传人的
//   2. pipeline 运行中产生的
//   3. pipeline 全局传人
//   4. action 默认默认值
func (e *Engine) run(req *runner.RunRequest, resp *runner.RunResponse) {
	if !e.init {
		resp.Failed("engine not init")
		return
	}

	s := req.Step

	e.log.Debugf("start run step: %s status %s", s.Key, s.Status)

	// 1.查询step对应的action定义
	descA := pipeline.NewDescribeActionRequest(s.GetNamespace(), s.ActionName(), s.ActionVersion())
	descA.Namespace = s.GetNamespace()
	ctx := gcontext.NewGrpcOutCtx()
	action, err := e.wc.Pipeline().DescribeAction(ctx.Context(), descA)
	if err != nil {
		resp.Failed("describe step action error, %s", err)
		return
	}

	// 2.加载Action默认参数
	req.LoadRunParams(action.DefaultRunParam())

	// 3.查询Pipeline, 加载全局参数
	if s.IsCreateByPipeline() {
		descP := pipeline.NewDescribePipelineRequestWithID(s.GetPipelineId())
		descP.Namespace = s.GetNamespace()
		pl, err := e.wc.Pipeline().DescribePipeline(ctx.Context(), descP)
		if err != nil {
			resp.Failed("describe step pipeline error, %s", err)
			return
		}
		req.LoadRunParams(pl.With)
		req.LoadMount(pl.Mount)
	}

	// 4. 加载step传递的参数
	req.LoadRunParams(s.With)

	// 校验run参数合法性
	if err := action.ValidateRunParam(req.RunParams); err != nil {
		resp.Failed(err.Error())
		return
	}

	// 加载Runner运行需要的参数
	req.LoadRunnerParams(action.RunnerParam())

	e.log.Debugf("choice %s runner to run step", action.RunnerType)
	// 3.根据action定义的runner_type, 调用具体的runner
	switch action.RunnerType {
	case pipeline.RUNNER_TYPE_DOCKER:
		e.docker.Run(context.Background(), req, resp)
	case pipeline.RUNNER_TYPE_K8s:
		e.k8s.Run(context.Background(), req, resp)
	case pipeline.RUNNER_TYPE_LOCAL:
		e.local.Run(context.Background(), req, resp)
	default:
		resp.Failed("unknown runner type: %s", action.RunnerType)
		return
	}

}

// 如果step执行完成
func (e *Engine) updateStep(s *pipeline.Step) {
	e.log.Debugf("receive step %s update, status %s", s.Key, s.Status)
	if err := e.recorder.Update(s.Clone()); err != nil {
		e.log.Errorf("update step status error, %s", err)
	}
}

package engine

import (
	"context"
	"fmt"

	"github.com/infraboard/mcube/grpc/gcontext"
	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/logger/zap"

	"github.com/infraboard/workflow/api/client"
	"github.com/infraboard/workflow/api/pkg/pipeline"
	"github.com/infraboard/workflow/common/informers/step"
	"github.com/infraboard/workflow/node/controller/step/runner"
	"github.com/infraboard/workflow/node/controller/step/runner/docker"
	"github.com/infraboard/workflow/node/controller/step/runner/k8s"
	"github.com/infraboard/workflow/node/controller/step/runner/local"
)

var (
	engine = &Engine{}
)

func RunStep(s *pipeline.Step) {
	engine.Run(s)
}

func Init(wc *client.Client, recorder step.Recorder) (err error) {
	if wc == nil {
		return fmt.Errorf("init runner error, workflow client is nil")
	}

	engine.recorder = recorder
	engine.wc = wc
	engine.docker, err = docker.NewRunner()
	engine.k8s = k8s.NewRunner()
	engine.local = local.NewRunner()
	engine.log = zap.L().Named("Runner.Engine")
	if err != nil {
		return err
	}

	engine.init = true
	return nil
}

type Engine struct {
	recorder step.Recorder
	wc       *client.Client
	docker   runner.Runner
	k8s      runner.Runner
	local    runner.Runner
	init     bool
	log      logger.Logger
}

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
		e.docker.Run(context.Background(), req)
	case pipeline.RUNNER_TYPE_K8s:
		e.k8s.Run(context.Background(), req)
	case pipeline.RUNNER_TYPE_LOCAL:
		e.local.Run(context.Background(), req)
	default:
		s.Failed("unknown runner type: %s", action.RunnerType)
		return
	}
}

func (e *Engine) updateStep(s *pipeline.Step) {
	if err := e.recorder.Update(s); err != nil {
		e.log.Errorf("update step status error, %s", err)
	}
}

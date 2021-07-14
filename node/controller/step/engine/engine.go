package engine

import (
	"context"
	"fmt"

	"github.com/infraboard/mcube/grpc/gcontext"

	"github.com/infraboard/workflow/api/client"
	"github.com/infraboard/workflow/api/pkg/pipeline"
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

func Init(wc *client.Client) (err error) {
	if wc == nil {
		return fmt.Errorf("init runner error, workflow client is nil")
	}

	engine.wc = wc
	engine.docker, err = docker.NewRunner()
	engine.k8s = k8s.NewRunner()
	engine.local = local.NewRunner()
	if err != nil {
		return err
	}

	engine.init = true
	return nil
}

type Engine struct {
	wc     *client.Client
	docker runner.Runner
	k8s    runner.Runner
	local  runner.Runner
	init   bool
}

// Run 运行Step
// step的参数加载优先级:
//   1. step 本身传人的
//   2. pipeline 运行中产生的
//   3. pipeline 全局传人
//   4. action 默认默认值
func (e *Engine) Run(s *pipeline.Step) {
	if !e.init {
		s.Failed(fmt.Errorf("engine not init"))
		return
	}

	// 构造运行请求
	req := runner.NewRunRequest(s)

	// 1.查询step对应的action定义
	descAction := pipeline.NewDescribeActionRequestWithName(s.Action)
	ctx := gcontext.NewGrpcOutCtx()
	action, err := e.wc.Pipeline().DescribeAction(ctx.Context(), descAction)
	if err != nil {
		s.Failed(err)
		return
	}

	// 加载运行时参数
	req.LoadRunParams(action.DefaultRunParam())

	// 校验参数合法性
	if err := action.ValidateParam(req.RunParams); err != nil {
		s.Failed(err)
		return
	}

	// 加载Runner运行需要的参数
	req.LoadRunnerParams(action.DefaultRunParam())

	// 3.根据action定义的runner_type, 调用具体的runner
	switch action.RunnerType {
	case pipeline.RUNNER_TYPE_DOCKER:
		err = e.docker.Run(context.Background(), req)
	case pipeline.RUNNER_TYPE_K8s:
		err = e.k8s.Run(context.Background(), req)
	case pipeline.RUNNER_TYPE_LOCAL:
		err = e.local.Run(context.Background(), req)
	default:
		s.Failed(fmt.Errorf("unknown runner type: %s", action.RunnerType))
		return
	}

	if err != nil {
		s.Failed(err)
		return
	}
}

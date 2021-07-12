package engine

import (
	"context"
	"fmt"

	"github.com/infraboard/mcube/exception"
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

func RunStep(s *pipeline.Step) error {
	return engine.Run(s)
}

func Init(wc *client.Client) error {
	if wc == nil {
		return fmt.Errorf("init runner error, workflow client is nil")
	}

	engine.wc = wc
	engine.docker = docker.NewRunner()
	engine.k8s = k8s.NewRunner()
	engine.local = local.NewRunner()
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

func (e *Engine) Run(s *pipeline.Step) error {
	if !e.init {
		return fmt.Errorf("runner engine not init")
	}

	fmt.Println(e.wc.Pipeline())

	// 1.查询step对应的action定义
	req := pipeline.NewDescribeActionRequestWithName(s.Name)
	ctx := gcontext.NewGrpcOutCtx()
	a, err := e.wc.Pipeline().DescribeAction(ctx.Context(), req)
	if err != nil {
		return err
	}

	// 2.根据action定义的runner_type, 调用具体的runner
	switch a.RunnerType {
	case pipeline.RUNNER_TYPE_DOCKER:
		return e.docker.Run(context.Background(), runner.NewRunRequest(s))
	case pipeline.RUNNER_TYPE_K8s:
		return e.k8s.Run(context.Background(), runner.NewRunRequest(s))
	case pipeline.RUNNER_TYPE_LOCAL:
		return e.local.Run(context.Background(), runner.NewRunRequest(s))
	default:
		return exception.NewBadRequest("unknown runner type: %s", a.RunnerType)
	}
}

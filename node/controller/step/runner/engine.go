package runner

import (
	"fmt"

	"github.com/infraboard/mcube/exception"

	"github.com/infraboard/workflow/api/client"
	"github.com/infraboard/workflow/api/pkg/pipeline"
	"github.com/infraboard/workflow/node/controller/step/runner/docker"
	"github.com/infraboard/workflow/node/controller/step/runner/k8s"
	"github.com/infraboard/workflow/node/controller/step/runner/local"
)

type Runner interface {
	Run(*pipeline.Step) error
}

var (
	engine = &Engine{}
)

func Init() error {
	engine.docker = docker.NewRunner()
	engine.k8s = k8s.NewRunner()
	engine.local = local.NewRunner()
	engine.init = true
	return nil
}

func RunStep(s *pipeline.Step) error {
	return engine.Run(s)
}

type Engine struct {
	wc     *client.Client
	docker Runner
	k8s    Runner
	local  Runner
	init   bool
}

func (e *Engine) Run(s *pipeline.Step) error {
	if !e.init {
		return fmt.Errorf("runner engine not init")
	}

	// 1.查询step对应的action定义
	a, err := e.wc.Pipeline().DescribeAction(nil, nil)
	if err != nil {
		return err
	}

	// 2.根据action定义的runner_type, 调用具体的runner
	switch a.RunnerType {
	case pipeline.RUNNER_TYPE_DOCKER:
		return e.docker.Run(s)
	case pipeline.RUNNER_TYPE_K8s:
		return e.k8s.Run(s)
	case pipeline.RUNNER_TYPE_LOCAL:
		return e.local.Run(s)
	default:
		return exception.NewBadRequest("unknown runner type: %s", a.RunnerType)
	}
}

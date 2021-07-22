package engine

import (
	"fmt"

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

func CancelStep(s *pipeline.Step) {

}

func Init(wc *client.Client, recorder step.Recorder) (err error) {
	if wc == nil {
		return fmt.Errorf("init runner error, workflow client is nil")
	}

	engine.log = zap.L().Named("Runner.Engine")
	engine.recorder = recorder
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
	recorder step.Recorder
	wc       *client.Client
	docker   runner.Runner
	k8s      runner.Runner
	local    runner.Runner
	init     bool
	log      logger.Logger
}

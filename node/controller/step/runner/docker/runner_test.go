package docker_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/infraboard/mcube/logger/zap"
	"github.com/infraboard/workflow/api/pkg/pipeline"
	"github.com/infraboard/workflow/node/controller/step/runner"
	"github.com/infraboard/workflow/node/controller/step/runner/docker"
)

var (
	dr *docker.Runner
)

var (
	smapleStep   = &pipeline.Step{Key: "test"}
	runnerParams = map[string]string{
		"IMAGE_URL": "busybox",
		"IMAGE_CMD": "date",
	}
	resp = runner.NewRunReponse(testUpdater)
)

func testUpdater(s *pipeline.Step) {
	fmt.Println(s)
}

func TestDockerRunNilStep(t *testing.T) {
	req := runner.NewRunRequest(nil)

	dr.Run(context.Background(), req, resp)
}

func TestDockerRunNULLStep(t *testing.T) {
	req := runner.NewRunRequest(&pipeline.Step{})
	dr.Run(context.Background(), req, resp)
	t.Log(req.Step)
}

func TestDockerRunSampleStep(t *testing.T) {
	req := runner.NewRunRequest(smapleStep)
	dr.Run(context.Background(), req, resp)
	t.Log(smapleStep)
}

func TestDockerRunStepWithRunnerParams(t *testing.T) {
	req := runner.NewRunRequest(smapleStep)
	req.LoadRunnerParams(runnerParams)
	dr.Run(context.Background(), req, resp)
	t.Log(smapleStep)
}

func init() {
	if err := zap.DevelopmentSetup(); err != nil {
		panic(err)
	}
	r, err := docker.NewRunner()
	if err != nil {
		panic(err)
	}
	dr = r

}

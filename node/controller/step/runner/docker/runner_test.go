package docker_test

import (
	"context"
	"fmt"
	"testing"
	"time"

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

func TestRunNilStep(t *testing.T) {
	req := runner.NewRunRequest(nil)

	dr.Run(context.Background(), req, resp)
}

func TestRunNULLStep(t *testing.T) {
	req := runner.NewRunRequest(&pipeline.Step{})
	dr.Run(context.Background(), req, resp)
	t.Log(req.Step)
}

func TestDockerRunSampleStep(t *testing.T) {
	req := runner.NewRunRequest(smapleStep)
	dr.Run(context.Background(), req, resp)
	t.Log(smapleStep)
}

func TestRunStepWithRunnerParams(t *testing.T) {
	req := runner.NewRunRequest(smapleStep)
	req.LoadRunnerParams(runnerParams)
	dr.Run(context.Background(), req, resp)
	t.Log(smapleStep)
}

func TestCancelStep(t *testing.T) {
	req := runner.NewRunRequest(smapleStep)
	req.LoadRunnerParams(cmdRunnerParams("busybox", "/bin/sleep,10"))
	go dr.Run(context.Background(), req, resp)

	time.Sleep(3 * time.Second)
	dr.Cancel(context.Background(), runner.NewCancelRequest(req.Step))
	// 等待容器退出
	time.Sleep(3 * time.Second)
	t.Log(resp)

}

func cmdRunnerParams(image, cmd string) map[string]string {
	return map[string]string{
		"IMAGE_URL": image,
		"IMAGE_CMD": cmd,
	}
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

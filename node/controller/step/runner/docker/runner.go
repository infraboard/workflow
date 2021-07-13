package docker

import (
	"context"
	"io"

	"github.com/docker/docker/client"
	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/logger/zap"

	"github.com/infraboard/workflow/node/controller/step/runner"
)

func NewRunner() (*Runner, error) {
	log := zap.L().Named("Runner.Docker")

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}

	log.Infof("docker runner connect success, version: %s", cli.ClientVersion())

	return &Runner{
		log: log,
		cli: cli,
	}, nil
}

type Runner struct {
	cli *client.Client
	log logger.Logger
}

func (r *Runner) Run(context.Context, *runner.RunRequest) error {
	return nil
}

func (r *Runner) Log(context.Context, *runner.LogRequest) (io.ReadCloser, error) {
	return nil, nil
}

func (r *Runner) Connect(context.Context, *runner.ConnectRequest) error {
	return nil
}

func (r *Runner) Cancel(context.Context, *runner.CancelRequest) error {
	return nil
}

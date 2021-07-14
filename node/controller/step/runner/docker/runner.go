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

// Runner Params:
//   IMAGE_URL: 镜像URL, 比如: docker-build
//   IMAGE_PULL_SECRET: 拉起镜像的凭证
//   IMAGE_PUSH_SECRET: 推送镜像的凭证
// Run Params:
//   IMAGE_VERSION: 镜像版本 比如: v1
//   GIT_SSH_URL: 代码仓库地址, 比如: git@gitee.com:infraboard/keyauth.git
//   IMAGE_PUSH_URL:
func (r *Runner) Run(ctx context.Context, req *runner.RunRequest) error {
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

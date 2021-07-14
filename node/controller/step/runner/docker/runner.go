package docker

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/logger/zap"

	"github.com/infraboard/workflow/api/pkg/pipeline"
	"github.com/infraboard/workflow/node/controller/step/runner"
)

const (
	IMAGE_URL_KEY     = "IMAGE_URL"
	IMAGE_VERSION_KEY = "IMAGE_VERSION"
)

const (
	CONTAINER_ID_KEY   = "container_id"
	CONTAINER_WARN_KEY = "container_warn"
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

// Docker官方SDK使用说明: https://docs.docker.com/engine/api/sdk/examples/
// Docker官方API使用说明: https://docs.docker.com/engine/api/v1.41/
type Runner struct {
	cli *client.Client
	log logger.Logger
}

// ContainerCreate参数说明:  https://docs.docker.com/engine/api/v1.41/#operation/ContainerCreate
// Runner Params:
//   IMAGE_URL: 镜像URL, 比如: docker-build
//   IMAGE_PULL_SECRET: 拉起镜像的凭证
//   IMAGE_PUSH_SECRET: 推送镜像的凭证
// Run Params:
//   IMAGE_VERSION: 镜像版本 比如: v1
//   GIT_SSH_URL: 代码仓库地址, 比如: git@gitee.com:infraboard/keyauth.git
//   IMAGE_PUSH_URL: 代码推送地址
func (r *Runner) Run(ctx context.Context, in *runner.RunRequest) {
	if in.Step == nil {
		r.log.Errorf("step is nil")
		return
	}

	req := newDockerRunRequest(in)
	if err := req.Validate(); err != nil {
		in.Step.Failed("validate docker run request error, %s", err)
		return
	}

	resp, err := r.cli.ContainerCreate(ctx, &container.Config{
		Image: req.Image(),
		Env:   req.ContainerEnv(),
	}, nil, nil, nil, req.ContainerName())

	if err != nil {
		in.Step.Failed("create container error, %s", err)
		return
	}

	err = r.cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{})
	if err != nil {
		in.Step.Failed("run container error, %s", err)
		return
	}

	in.Step.Success(map[string]string{
		CONTAINER_ID_KEY:   resp.ID,
		CONTAINER_WARN_KEY: strings.Join(resp.Warnings, ","),
	})
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

func newDockerRunRequest(r *runner.RunRequest) *dockerRunRequest {
	if r.Step.Status == nil {
		r.Step.Status = pipeline.NewDefaultStepStatus()
	}
	return &dockerRunRequest{r}
}

type dockerRunRequest struct {
	*runner.RunRequest
}

func (r *dockerRunRequest) Image() string {
	if r.ImageVersion() == "" {
		return r.ImageURL()
	}
	return fmt.Sprintf("%s:%s", r.ImageURL(), r.ImageVersion())
}

func (r *dockerRunRequest) ImageURL() string {
	return r.RunnerParams[IMAGE_URL_KEY]
}

func (r *dockerRunRequest) ImageVersion() string {
	return r.RunParams[IMAGE_VERSION_KEY]
}

func (r *dockerRunRequest) ContainerName() string {
	return r.Step.Key
}

func (r *dockerRunRequest) ContainerEnv() []string {
	envs := []string{}
	for k, v := range r.mergeParams() {
		envs = append(envs, fmt.Sprintf("%s=%s", k, v))
	}
	return envs
}

func (r *dockerRunRequest) mergeParams() map[string]string {
	m := r.RunParams
	for k, v := range r.Step.With {
		m[k] = v
	}
	return m
}

func (r *dockerRunRequest) Validate() error {
	fmt.Println(r.Image())
	if r.ImageURL() == "" {
		return fmt.Errorf("%s missed", IMAGE_URL_KEY)
	}
	return nil
}

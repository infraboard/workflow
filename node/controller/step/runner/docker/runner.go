package docker

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/logger/zap"

	"github.com/infraboard/workflow/api/apps/action"
	"github.com/infraboard/workflow/node/controller/step/runner"
	"github.com/infraboard/workflow/node/controller/step/store"
)

const (
	IMAGE_URL_KEY     = "IMAGE_URL"
	IMAGE_CMD_KEY     = "IMAGE_CMD"
	IMAGE_VERSION_KEY = "IMAGE_VERSION"
)

var (
	IMAGE_URL_KEY_DESC = &action.RunParamDesc{
		KeyName:   IMAGE_URL_KEY,
		KeyDesc:   "镜像地址",
		Required:  true,
		ValueDesc: "镜像仓库地址, 比如busybox, 如果是私有仓库请配置认证信息",
	}
	IMAGE_VERSION_KEY_DESC = &action.RunParamDesc{
		KeyName:   IMAGE_VERSION_KEY,
		KeyDesc:   "镜像版本",
		Required:  false,
		ValueDesc: "镜像对应的Tag, 比如latest",
	}
	IMAGE_CMD_KEY_DESC = &action.RunParamDesc{
		KeyName:   IMAGE_CMD_KEY,
		KeyDesc:   "执行命令",
		Required:  true,
		ValueDesc: "如果是多部分请用逗号分隔, 比如 sleep,10",
	}
)

func ParamsDesc() []*action.RunParamDesc {
	return []*action.RunParamDesc{
		IMAGE_URL_KEY_DESC,
		IMAGE_VERSION_KEY_DESC,
		IMAGE_CMD_KEY_DESC,
	}
}

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
	ctm := 3 * time.Second

	return &Runner{
		log:           log,
		cli:           cli,
		store:         store.NewStore(),
		cancelTimeout: &ctm,
	}, nil
}

// Docker官方SDK使用说明: https://docs.docker.com/engine/api/sdk/examples/
// Docker官方API使用说明: https://docs.docker.com/engine/api/v1.41/
type Runner struct {
	cli           *client.Client
	log           logger.Logger
	store         store.StoreFactory
	cancelTimeout *time.Duration
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
func (r *Runner) Run(ctx context.Context, in *runner.RunRequest, out *runner.RunResponse) {
	r.log.Debugf("docker start run step: %s", in.Step.Key)

	req := newDockerRunRequest(in)
	if err := req.Validate(); err != nil {
		out.Failed("validate docker run request error, %s", err)
		return
	}

	// 创建容器
	resp, err := r.cli.ContainerCreate(ctx, &container.Config{
		Image: req.Image(),
		Env:   req.ContainerEnv(),
		Cmd:   req.ContainerCMD(),
	}, nil, nil, nil, req.ContainerName())
	if err != nil {
		out.Failed("create container error, %s", err)
		return
	}

	// 退出时销毁容器
	defer r.removeContainer(resp.ID)

	// 更新状态
	up := r.store.NewFileUploader(req.Step.Key)
	out.UpdateReponseMap("log_driver", up.DriverName())
	out.UpdateReponseMap("log_path", up.ObjectID())
	out.UpdateReponseMap(CONTAINER_ID_KEY, resp.ID)
	out.UpdateReponseMap(CONTAINER_WARN_KEY, strings.Join(resp.Warnings, ","))
	out.UpdateResponse(in.Step)

	// TODO: 如果镜像不存在, 要提前拉去

	// 启动容器
	err = r.cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{})
	if err != nil {
		// 启动失败则删除容器
		r.removeContainer(resp.ID)
		out.Failed("run container error, %s", err)
		return
	}

	// 记录容器的日志
	logStream, err := r.cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true, ShowStderr: true})
	if err != nil {
		out.Failed("get container log error, %s", err)
		return
	}

	// 上传容器日志
	if err := up.Upload(ctx, logStream); err != nil {
		out.Failed(err.Error())
		return
	}

	// 等待容器退出
	statusCh, errCh := r.cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			out.Failed(err.Error())
			return
		}
	case <-statusCh:
	}

	// 容器退出
	if err := r.containerExit(resp.ID); err != nil {
		out.Failed(err.Error())
		return
	}
}

// 容器退出时, 需要
// 1. 判断容器执行成功还是失败
// 2. 收集容器运行时产生的日志
// 3. 收集容器执行时的输出结果
func (r *Runner) containerExit(id string) error {
	r.log.Debugf("container %s strart exit ...", id)

	info, err := r.inspectContainer(id)
	if err != nil {
		return fmt.Errorf("inspec container error, %s", err)
	}

	// 获取容器执行退出状态
	state := info.State
	if state.ExitCode != 0 {
		msg := fmt.Sprintf("container run failed, status %s, exit code is %d", state.Status, state.ExitCode)
		if info.State.Error != "" {
			msg += fmt.Sprintf(", error: %s", state.Error)
		}
		return fmt.Errorf(msg)
	}

	// 通过挂入的卷 收集容器的返回
	return nil
}

func (r *Runner) inspectContainer(id string) (*types.ContainerJSON, error) {
	resp, err := r.cli.ContainerInspect(context.Background(), id)
	return &resp, err
}

// 删除容器
func (r *Runner) removeContainer(id string) {
	err := r.cli.ContainerRemove(context.Background(), id, types.ContainerRemoveOptions{})
	if err != nil {
		r.log.Errorf("remove contain %s failed", err)
	}
}

func (r *Runner) Cancel(ctx context.Context, in *runner.CancelRequest) {
	req := newDockerCancelRequest(in)
	if err := req.Validate(); err != nil {
		in.Step.Failed("validate container cancel request error, %s", err)
		return
	}

	if err := r.cli.ContainerStop(ctx, req.ContainerID(), r.cancelTimeout); err != nil {
		in.Step.Failed("cancel container error, %s", err)
		return
	}
}

func (r *Runner) Log(context.Context, *runner.LogRequest) (io.ReadCloser, error) {
	return nil, nil
}

func (r *Runner) Connect(context.Context, *runner.ConnectRequest) error {
	return nil
}

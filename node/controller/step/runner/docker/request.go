package docker

import (
	"fmt"
	"strings"

	"github.com/infraboard/workflow/api/app/pipeline"
	"github.com/infraboard/workflow/node/controller/step/runner"
)

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
	return r.RunnerParams[IMAGE_VERSION_KEY]
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

func (r *dockerRunRequest) ContainerCMD() []string {
	return strings.Split(r.RunnerParams[IMAGE_CMD_KEY], ",")
}

func (r *dockerRunRequest) mergeParams() map[string]string {
	m := r.RunParams
	for k, v := range r.Step.With {
		m[k] = v
	}
	return m
}

func (r *dockerRunRequest) Validate() error {
	if r.Step == nil || r.Step.Key == "" {
		return fmt.Errorf("step is nil or step key is \"\"")
	}

	if r.ImageURL() == "" {
		return fmt.Errorf("%s missed", IMAGE_URL_KEY)
	}

	return nil
}

func newDockerCancelRequest(r *runner.CancelRequest) *dockerCancelRequest {
	if r.Step.Status == nil {
		r.Step.Status = pipeline.NewDefaultStepStatus()
	}
	return &dockerCancelRequest{r}
}

type dockerCancelRequest struct {
	*runner.CancelRequest
}

func (r *dockerCancelRequest) ContainerID() string {
	if r.Step == nil || r.Step.Status == nil || r.Step.Status.Response == nil {
		return ""
	}

	return r.Step.Status.Response[CONTAINER_ID_KEY]
}

func (r *dockerCancelRequest) Validate() error {
	if r.Step == nil || r.Step.Key == "" {
		return fmt.Errorf("step is nil or step key is \"\"")
	}

	if r.ContainerID() == "" {
		return fmt.Errorf("%s missed", CONTAINER_ID_KEY)
	}

	return nil
}

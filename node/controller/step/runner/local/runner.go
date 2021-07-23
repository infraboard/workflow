package local

import (
	"context"
	"io"

	"github.com/infraboard/workflow/node/controller/step/runner"
)

func NewRunner() *Runner {
	return &Runner{}
}

type Runner struct {
}

func (r *Runner) Run(ctx context.Context, in *runner.RunRequest, out *runner.RunResponse) {
}

func (r *Runner) Log(context.Context, *runner.LogRequest) (io.ReadCloser, error) {
	return nil, nil
}

func (r *Runner) Connect(context.Context, *runner.ConnectRequest) error {
	return nil
}

func (r *Runner) Cancel(context.Context, *runner.CancelRequest) {
	return
}

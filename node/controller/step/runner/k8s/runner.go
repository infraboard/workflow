package k8s

import "github.com/infraboard/workflow/api/pkg/pipeline"

func NewRunner() *Runner {
	return &Runner{}
}

type Runner struct{}

func (r *Runner) Run(s *pipeline.Step) error {
	return nil
}

package runner

import "github.com/infraboard/workflow/api/pkg/pipeline"

type Runner interface {
	Run(*pipeline.Step) error
}

var (
	engine = &Engine{}
)

func RunStep(s *pipeline.Step) error {
	return engine.Run(s)
}

type Engine struct{}

func (e *Engine) Run(s *pipeline.Step) error {
	// 1.查询step 对应的action定义

	// 2.根据action定义的runner_type, 调用具体的runner

	return nil
}

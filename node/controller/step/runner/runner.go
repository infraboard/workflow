package runner

import (
	"context"
	"io"

	"github.com/infraboard/workflow/api/pkg/pipeline"
)

type Runner interface {
	// 执行Step, 执行过后的关联信息保存在Status的Response里面
	Run(context.Context, *RunRequest)
	// 连接到该执行环境
	Connect(context.Context, *ConnectRequest) error
	// 取消Step的执行
	Cancel(context.Context, *CancelRequest) error
}

func NewRunRequest(s *pipeline.Step) *RunRequest {
	return &RunRequest{
		Step:         s,
		RunnerParams: map[string]string{},
		RunParams:    map[string]string{},
	}
}

type RunRequest struct {
	RunnerParams map[string]string   // runner 运行需要的参数
	RunParams    map[string]string   // Pipeline定义的全局变量
	Mount        *pipeline.MountData // pipeline定义的挂载文件
	Step         *pipeline.Step      // 具体step
}

func (r *RunRequest) LoadMount(m *pipeline.MountData) {
	r.Mount = m
}

func (r *RunRequest) LoadRunParams(params map[string]string) {
	for k, v := range params {
		r.RunParams[k] = v
	}
}

func (r *RunRequest) LoadRunnerParams(params map[string]string) {
	for k, v := range params {
		r.RunnerParams[k] = v
	}
}

type LogRequest struct {
	Step *pipeline.Step
}

type CancelRequest struct {
	Step *pipeline.Step
}

// // ConnectRequest holds information pertaining to the current streaming session:
// // input/output streams, if the client is requesting a TTY, and a terminal size queue to
// // support terminal resizing.
type ConnectRequest struct {
	Step              *pipeline.Step
	Stdin             io.Reader
	Stdout            io.Writer
	Stderr            io.Writer
	Tty               bool
	TerminalSizeQueue TerminalSizeQueue
}

// TerminalSize and TerminalSizeQueue was a part of k8s.io/kubernetes/pkg/util/term
// and were moved in order to decouple client from other term dependencies

// TerminalSize represents the width and height of a terminal.
type TerminalSize struct {
	Width  uint16
	Height uint16
}

// TerminalSizeQueue is capable of returning terminal resize events as they occur.
type TerminalSizeQueue interface {
	// Next returns the new terminal size after the terminal has been resized. It returns nil when
	// monitoring has been stopped.
	Next() *TerminalSize
}

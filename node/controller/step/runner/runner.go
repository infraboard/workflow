package runner

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/infraboard/workflow/api/pkg/pipeline"
)

type Runner interface {
	// 执行Step, 执行过后的关联信息保存在Status的Response里面
	Run(context.Context, *RunRequest, *RunResponse)
	// 连接到该执行环境
	Connect(context.Context, *ConnectRequest) error
	// 取消Step的执行
	Cancel(context.Context, *CancelRequest)
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
	RunParams    map[string]string   // step 运行需要的参数
	Mount        *pipeline.MountData // 挂载文件
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

func NewRunReponse(updater UpdateStepCallback) *RunResponse {
	return &RunResponse{
		updater: updater,
		resp:    map[string]string{},
		ctx:     map[string]string{},
	}
}

type UpdateStepCallback func(*pipeline.Step)

type RunResponse struct {
	updater UpdateStepCallback // 更新状态的回调
	errs    []string
	resp    map[string]string
	ctx     map[string]string
}

func (r *RunResponse) UpdateReponseMap(k, v string) {
	r.resp[k] = v
}

func (r *RunResponse) UpdateCtxMap(k, v string) {
	r.ctx[k] = v
}

func (r *RunResponse) UpdateResponse(s *pipeline.Step) {
	r.updater(s)
}

func (r *RunResponse) Failed(format string, a ...interface{}) {
	r.errs = append(r.errs, fmt.Sprintf(format, a...))
}

func (r *RunResponse) HasError() bool {
	return len(r.errs) > 0
}

func (r *RunResponse) ErrorMessage() string {
	return strings.Join(r.errs, ",")
}

type LogRequest struct {
	Step *pipeline.Step
}

func NewCancelRequest(s *pipeline.Step) *CancelRequest {
	return &CancelRequest{
		Step: s,
	}
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

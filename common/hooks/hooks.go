package hooks

import (
	"context"

	"github.com/infraboard/workflow/api/pkg/pipeline"
	"github.com/infraboard/workflow/common/hooks/webhook"
)

type StepNotifyEvent struct {
	StepKey      string            `json:"step_key"`
	NotifyParams map[string]string `json:"notify_params"`
	*pipeline.StepStatus
}

// StepWebHooker step状态变化时，通知其他系统
type StepWebHookPusher interface {
	Send(context.Context, []*pipeline.WebHook, *pipeline.Step) error
}

func NewDefaultStepWebHookPusher() StepWebHookPusher {
	return &webhook.WebHook{}
}

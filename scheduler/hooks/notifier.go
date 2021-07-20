package notify

import (
	"context"

	"github.com/infraboard/workflow/api/pkg/pipeline"
)

type StepNotifyEvent struct {
	StepKey      string            `json:"step_key"`
	NotifyParams map[string]string `json:"notify_params"`
	*pipeline.StepStatus
}

// Notifier 用于将
type Notifier interface {
	Send(context.Context, *pipeline.StepStatus) error
}

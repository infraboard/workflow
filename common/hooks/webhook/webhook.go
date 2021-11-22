package webhook

import (
	"context"
	"fmt"

	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/logger/zap"

	"github.com/infraboard/workflow/api/app/pipeline"
)

func NewWebHook() *WebHook {
	return &WebHook{
		log: zap.L().Named("WebHook"),
	}
}

type WebHook struct {
	log logger.Logger
}

func (h *WebHook) Send(ctx context.Context, hooks []*pipeline.WebHook, step *pipeline.Step) error {
	if step == nil {
		return fmt.Errorf("step is nil")
	}

	if err := h.validate(hooks); err != nil {
		return err
	}

	h.log.Debugf("start send step[%s] webhook, total %d", step.Key, len(hooks))
	for i := range hooks {
		req := newRequest(hooks[i], step)
		req.Push()
	}

	return nil
}

func (h *WebHook) validate(hooks []*pipeline.WebHook) error {
	if len(hooks) == 0 {
		return nil
	}

	if len(hooks) > MAX_WEBHOOKS_PER_SEND {
		return fmt.Errorf("too many webhooks configs current: %d, max: %d", len(hooks), MAX_WEBHOOKS_PER_SEND)
	}

	return nil
}

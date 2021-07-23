package webhook

import (
	"context"
	"fmt"

	"github.com/infraboard/workflow/api/pkg/pipeline"
)

type WebHook struct{}

func (h *WebHook) Send(ctx context.Context, hooks []*pipeline.WebHook, step *pipeline.Step) error {
	if step == nil {
		return fmt.Errorf("step is nil")
	}

	if err := h.validate(hooks); err != nil {
		return err
	}

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

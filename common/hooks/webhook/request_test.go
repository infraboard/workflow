package webhook_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/infraboard/workflow/api/pkg/pipeline"
	"github.com/infraboard/workflow/common/hooks/webhook"
)

var (
	botURL = "https://open.feishu.cn/open-apis/bot/v2/hook/xxx"
)

func TestFeishuWebHook(t *testing.T) {
	should := assert.New(t)

	hooks := testPipelineWebHook()
	sender := webhook.NewWebHook()
	err := sender.Send(
		context.Background(),
		hooks,
		testPipelineStep(),
	)
	should.NoError(err)

	t.Log(hooks[0])
}

func testPipelineWebHook() []*pipeline.WebHook {
	h1 := &pipeline.WebHook{
		Url:         botURL,
		Events:      []pipeline.STEP_STATUS{pipeline.STEP_STATUS_SUCCEEDED},
		Description: "测试",
	}
	return []*pipeline.WebHook{h1}
}

func testPipelineStep() *pipeline.Step {
	return &pipeline.Step{
		Name: "only for test",
		Status: &pipeline.StepStatus{
			Status: pipeline.STEP_STATUS_SUCCEEDED,
		},
	}
}

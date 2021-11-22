package webhook_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/infraboard/workflow/api/app/pipeline"
	"github.com/infraboard/workflow/common/hooks/webhook"
)

var (
	feishuBotURL = "https://open.feishu.cn/open-apis/bot/v2/hook/83bde95c-00b2-4df1-91e4-705f66102479"
	dingBotURL   = "https://oapi.dingtalk.com/robot/send?access_token=xxxx"
	wechatBotURL = "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=693axxx6-7aoc-4bc4-97a0-0ec2sifa5aaa"
)

func TestFeishuWebHook(t *testing.T) {
	should := assert.New(t)

	hooks := testPipelineWebHook(feishuBotURL)
	sender := webhook.NewWebHook()
	err := sender.Send(
		context.Background(),
		hooks,
		testPipelineStep(),
	)
	should.NoError(err)

	t.Log(hooks[0])
}

func TestDingDingWebHook(t *testing.T) {
	should := assert.New(t)

	hooks := testPipelineWebHook(dingBotURL)
	sender := webhook.NewWebHook()
	err := sender.Send(
		context.Background(),
		hooks,
		testPipelineStep(),
	)
	should.NoError(err)

	t.Log(hooks[0])
}

func TestWechatWebHook(t *testing.T) {
	should := assert.New(t)

	hooks := testPipelineWebHook(wechatBotURL)
	sender := webhook.NewWebHook()
	err := sender.Send(
		context.Background(),
		hooks,
		testPipelineStep(),
	)
	should.NoError(err)
	t.Log(hooks[0])
}

func testPipelineWebHook(url string) []*pipeline.WebHook {
	h1 := &pipeline.WebHook{
		Url:         url,
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

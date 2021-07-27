package webhook

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/infraboard/workflow/api/pkg/pipeline"
	"github.com/infraboard/workflow/common/hooks/webhook/feishu"
)

const (
	MAX_WEBHOOKS_PER_SEND = 12
)

const (
	feishuBot = "feishu"
)

var (
	client = &http.Client{
		Timeout: 3 * time.Second,
	}
)

func newRequest(hook *pipeline.WebHook, step *pipeline.Step) *request {
	return &request{
		hook: hook,
		step: step,
	}
}

type request struct {
	hook *pipeline.WebHook
	step *pipeline.Step
}

func (r *request) Push() {
	r.hook.StartSend()

	// 准备请求,适配主流机器人
	var messageObj interface{}
	switch r.BotType() {
	case feishuBot:
		messageObj = feishu.NewStepCardMessage(r.step)
	default:
		messageObj = r.step
	}

	body, err := json.Marshal(messageObj)
	if err != nil {
		r.hook.SendFailed("marshal step to json error, %s", err)
		return
	}

	req, err := http.NewRequest("POST", r.hook.Url, bytes.NewReader(body))
	if err != nil {
		r.hook.SendFailed("new post request error, %s", err)
		return
	}

	for k, v := range r.hook.Header {
		req.Header.Add(k, v)
	}

	// 发起请求
	resp, err := client.Do(req)
	if err != nil {
		r.hook.SendFailed("send request error, %s", err)
		return
	}
	defer resp.Body.Close()

	if (resp.StatusCode / 100) != 2 {
		r.hook.SendFailed("status code[%d] is not 200", resp.StatusCode)
		return
	}

	r.hook.Success()
}

func (r *request) BotType() string {
	if strings.HasPrefix(r.hook.Url, feishu.URL_PREFIX) {
		return feishuBot
	}
	return ""
}

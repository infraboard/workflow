package feishu

import "github.com/infraboard/workflow/api/pkg/pipeline"

const (
	URL_PREFIX = "https://open.feishu.cn/open-apis/bot"
)

const (
	CardMessage = "interactive"
)

type MessageType string

func NewStepCardMessage(s *pipeline.Step) *Message {
	return &Message{
		MsgType: CardMessage,
		Card: Card{
			Config:   messageConfig(),
			Header:   messageHeader(s),
			Elements: messageContent(s),
		},
	}
}

// https://www.feishu.cn/hc/zh-CN/articles/360024984973
// 默认使用飞书的card数据模式
type Message struct {
	MsgType MessageType `json:"msg_type"`
	Card    Card        `json:"card"`
}

func messageContent(s *pipeline.Step) []*Element {
	content := NewMarkdownContent(s.String())
	return []*Element{content}
}

func messageHeader(s *pipeline.Step) *CardHeader {
	return &CardHeader{
		Title: map[string]string{
			"tag":     "plain_text",
			"content": s.ShowTitle(),
		},
		Template: messageCardColor(s),
	}
}

func messageConfig() *CardConfig {
	return &CardConfig{
		WideScreenMode: true,
		EnableForward:  true,
	}
}

func messageCardColor(s *pipeline.Step) string {
	if s.Status == nil {
		s.Status = pipeline.NewDefaultStepStatus()
	}

	switch s.Status.Status {
	case pipeline.STEP_STATUS_PENDDING,
		pipeline.STEP_STATUS_SKIP:
		return "grey"
	case pipeline.STEP_STATUS_RUNNING:
		return "turquoise"
	case pipeline.STEP_STATUS_SUCCEEDED:
		return "green"
	case pipeline.STEP_STATUS_FAILED,
		pipeline.STEP_STATUS_REFUSE:
		return "red"
	case pipeline.STEP_STATUS_CANCELING,
		pipeline.STEP_STATUS_CANCELED:
		return "yellow"
	case pipeline.STEP_STATUS_AUDITING:
		return "purple"
	default:
		return "wathet"
	}
}

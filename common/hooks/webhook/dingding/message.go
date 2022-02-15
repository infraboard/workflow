package dingding

import "github.com/infraboard/workflow/api/apps/pipeline"

const (
	URL_PREFIX = "https://oapi.dingtalk.com/robot/send"
)

const (
	CardMessage = "actionCard"
)

type MessageType string

func NewStepCardMessage(s *pipeline.Step) *Message {
	return &Message{
		MsgType:    CardMessage,
		ActionCard: newActionCard(s),
	}
}

// 自定义机器人接入: https://developers.dingtalk.com/document/app/custom-robot-access
// 默认使用钉钉的actionCard数据模式
type Message struct {
	MsgType    MessageType `json:"msgtype"`
	ActionCard *ActionCard `json:"actionCard"`
}

func newActionCard(s *pipeline.Step) *ActionCard {
	return &ActionCard{
		Title:              s.ShowTitle(),
		Text:               s.Status.String(),
		ButtonsOrientation: "0",
		SingleTitle:        "详情",
		SingleURL:          "https://www.dingtalk.com/",
	}
}

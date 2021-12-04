package feishu

import "fmt"

const (
	URL_PREFIX = "https://open.feishu.cn/open-apis/bot"
)
const (
	CardMessage = "interactive"
)

type MessageType string

func NewMarkdownNotifyMessage(robotURL, title, content string) *NotifyMessage {
	return &NotifyMessage{
		Title:    title,
		Content:  content,
		RobotURL: robotURL,
		Note:     []string{},
	}
}
func NewFiledMarkdownMessage(robotURL, title string, color Color, groups ...*FiledGroup) *NotifyMessage {
	return &NotifyMessage{
		Title:      title,
		FiledGroup: groups,
		RobotURL:   robotURL,
		Color:      color,
		Note:       []string{},
	}
}

type NotifyMessage struct {
	Title      string
	Content    string
	RobotURL   string
	FiledGroup []*FiledGroup
	Note       []string
	Color      Color
}

func (m *NotifyMessage) HasFiledGroup() bool {
	return len(m.FiledGroup) > 0
}

func (m *NotifyMessage) HasNote() bool {
	return len(m.Note) > 0
}

func (m *NotifyMessage) AddFiledGroup(group *FiledGroup) {
	m.FiledGroup = append(m.FiledGroup, group)
}

func (m *NotifyMessage) AddNote(n string) {
	m.Note = append(m.Note, n)
}

type FiledGroupEndType string

const (
	FiledGroupEndType_None FiledGroupEndType = "none"
	FiledGroupEndType_Hr   FiledGroupEndType = "hr"
	FiledGroupEndType_Line FiledGroupEndType = "line"
)

func NewEndHrGroup(fileds []*NotifyFiled) *FiledGroup {
	return &FiledGroup{
		EndType: FiledGroupEndType_Hr,
		Items:   fileds,
	}
}

func NewEndNoneGroup() *FiledGroup {
	return &FiledGroup{
		EndType: FiledGroupEndType_None,
		Items:   []*NotifyFiled{},
	}
}

type FiledGroup struct {
	EndType FiledGroupEndType
	Items   []*NotifyFiled
}

func (g *FiledGroup) Add(f *NotifyFiled) {
	g.Items = append(g.Items, f)
}

func NewNotifyFiled(key, value string, short bool) *NotifyFiled {
	return &NotifyFiled{
		Key:     key,
		Value:   value,
		IsShort: short,
	}
}

type NotifyFiled struct {
	IsShort bool
	Key     string
	Value   string
}

func (f *NotifyFiled) FiledFormat() string {
	return fmt.Sprintf("**%s**\n%s", f.Key, f.Value)
}

// 如何寻找emoji字符: https://emojipedia.org/light-bulb/
func NewCardMessage(m *NotifyMessage) *Message {
	return &Message{
		MsgType: CardMessage,
		Card: Card{
			Config:   messageConfig(),
			Header:   messageHeader(m),
			Elements: messageContent(m),
		},
	}
}

// https://www.feishu.cn/hc/zh-CN/articles/360024984973
// 默认使用飞书的card数据模式
type Message struct {
	MsgType MessageType `json:"msg_type"`
	Card    Card        `json:"card"`
}

// https://open.feishu.cn/document/ukTMukTMukTM/uEjNwUjLxYDM14SM2ATN
func messageContent(m *NotifyMessage) (elements []interface{}) {
	// 内容模块
	if m.HasFiledGroup() {
		for i := range m.FiledGroup {
			group := m.FiledGroup[i]
			content := NewFiledMarkdownContent(group.Items)
			elements = append(elements, content)

			switch group.EndType {
			case FiledGroupEndType_Hr:
				elements = append(elements, NewHrElement())
			case FiledGroupEndType_Line:
				content.Fields = append(content.Fields, NewField(false, ""))
			}
		}

	} else {
		content := NewMarkdownContent(m.Content)
		elements = append(elements, content)
	}
	// 备注模块
	if m.HasNote() {
		note := NewNoteContent(m.Note)
		elements = append(elements, NewHrElement(), note)
	}
	return
}

func messageHeader(m *NotifyMessage) *CardHeader {
	return &CardHeader{
		Title: map[string]string{
			"tag":     "plain_text",
			"content": m.Title,
		},
		Template: m.Color.String(),
	}
}

func messageConfig() *CardConfig {
	return &CardConfig{
		WideScreenMode: true,
		EnableForward:  true,
	}
}

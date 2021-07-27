package feishu

// card说明: https://open.feishu.cn/document/ukTMukTMukTM/ugTNwUjL4UDM14CO1ATN
// card可视化工具: https://open.feishu.cn/tool/cardbuilder?from=custom_bot_doc
type Card struct {
	Config   *CardConfig `json:"config"`
	Header   *CardHeader `json:"header"`
	Elements []*Element  `json:"elements"`
}

// card config说明: https://open.feishu.cn/document/ukTMukTMukTM/uAjNwUjLwYDM14CM2ATN
type CardConfig struct {
	WideScreenMode bool `json:"wide_screen_mode"`
	EnableForward  bool `json:"enable_forward"`
}

// card header说明: https://open.feishu.cn/document/ukTMukTMukTM/ukTNwUjL5UDM14SO1ATN
type CardHeader struct {
	Title    map[string]string `json:"title"`
	Template string            `json:"template"`
}

func NewMarkdownContent(content string) *Element {
	return &Element{
		Tag: "div",
		Text: Text{
			Content: content,
			Tag:     "lark_md",
		},
	}
}

// car element说明: https://open.feishu.cn/document/ukTMukTMukTM/uEjNwUjLxYDM14SM2ATN
type Element struct {
	Tag    string     `json:"tag"`
	Text   Text       `json:"text"`
	Fields FieldGroup `json:"fields"`
}

// 说明文档: https://open.feishu.cn/document/ukTMukTMukTM/uUzNwUjL1cDM14SN3ATN
type Text struct {
	Tag     string `json:"tag"`
	Content string `json:"content"`
	Lines   int    `json:"lines"`
}

// 说明文档 https://open.feishu.cn/document/ukTMukTMukTM/uYzNwUjL2cDM14iN3ATN
type FieldGroup struct {
	IsShort bool `json:"is_short"`
	Text    Text `json:"text"`
}

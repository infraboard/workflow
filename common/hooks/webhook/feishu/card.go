package feishu

const (
	MarkdownTagKey = "lark_md"
	DivTagKey      = "div"
)

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
		Tag: DivTagKey,
		Text: &Text{
			Content: content,
			Tag:     MarkdownTagKey,
		},
	}
}

func NewFieldsElement() *Element {
	return &Element{
		Tag:    DivTagKey,
		Fields: []*Field{},
	}
}

// car element说明: https://open.feishu.cn/document/ukTMukTMukTM/uEjNwUjLxYDM14SM2ATN
type Element struct {
	Tag    string   `json:"tag"`
	Text   *Text    `json:"text"`
	Fields []*Field `json:"fields"`
}

func (e *Element) AddField(f *Field) {
	e.Fields = append(e.Fields, f)
}

// 说明文档: https://open.feishu.cn/document/ukTMukTMukTM/uUzNwUjL1cDM14SN3ATN
type Text struct {
	Tag     string `json:"tag"`
	Content string `json:"content"`
	Lines   int    `json:"lines"`
}

func NewField(isShort bool, content string) *Field {
	return &Field{
		IsShort: isShort,
		Text: Text{
			Content: content,
			Tag:     MarkdownTagKey,
		},
	}
}

// 说明文档 https://open.feishu.cn/document/ukTMukTMukTM/uYzNwUjL2cDM14iN3ATN
type Field struct {
	IsShort bool `json:"is_short"`
	Text    Text `json:"text"`
}

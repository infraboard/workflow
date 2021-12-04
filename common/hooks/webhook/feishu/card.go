package feishu

// card说明: https://open.feishu.cn/document/ukTMukTMukTM/ugTNwUjL4UDM14CO1ATN
// card可视化工具: https://open.feishu.cn/tool/cardbuilder?from=custom_bot_doc
type Card struct {
	Config   *CardConfig   `json:"config"`
	Header   *CardHeader   `json:"header"`
	Elements []interface{} `json:"elements"`
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

type ElementType string

const (
	ElementType_Content = "content"
	ElementType_Hr      = "hr"
	ElementType_Image   = "image"
	ElementType_Action  = "action"
	ElementType_Note    = "note"
)

func NewMarkdownContent(content string) *ContentElement {
	return &ContentElement{
		Tag: "div",
		Text: &Text{
			Content: content,
			Tag:     "lark_md",
		},
	}
}
func NewFiledMarkdownContent(fileds []*NotifyFiled) *ContentElement {
	element := &ContentElement{
		Tag:    "div",
		Fields: []*Field{},
	}
	for i := range fileds {
		element.Fields = append(element.Fields, NewField(fileds[i].IsShort, fileds[i].FiledFormat()))
	}
	return element
}

// car element说明: https://open.feishu.cn/document/ukTMukTMukTM/uEjNwUjLxYDM14SM2ATN
type ContentElement struct {
	Tag    string   `json:"tag"`
	Text   *Text    `json:"text"`
	Fields []*Field `json:"fields"`
}

func NewNoteContent(fileds []string) *NoteElement {
	element := &NoteElement{
		Tag:      "note",
		Elements: []*Text{},
	}
	for i := range fileds {
		element.Elements = append(element.Elements, &Text{
			Content: fileds[i],
			Tag:     "plain_text",
		})
	}
	return element
}

// https://open.feishu.cn/document/ukTMukTMukTM/ucjNwUjL3YDM14yN2ATN
type NoteElement struct {
	Tag      string  `json:"tag"`
	Elements []*Text `json:"elements"`
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
			Tag:     "lark_md",
		},
	}
}

// 说明文档 https://open.feishu.cn/document/ukTMukTMukTM/uYzNwUjL2cDM14iN3ATN
type Field struct {
	IsShort bool `json:"is_short"`
	Text    Text `json:"text"`
}

func NewHrElement() *HrElement {
	return &HrElement{
		Tag: "hr",
	}
}

type HrElement struct {
	Tag string `json:"tag"`
}

package feishu

import (
	"fmt"

	"github.com/chyroc/lark"
)

func NewClient(appId, appSecret string) *Client {
	fmt.Println(appId)
	return &Client{
		client: lark.New(lark.WithAppCredential(appId, appSecret)),
	}
}

// 获取应用身份访问凭证 https://open.feishu.cn/document/ukTMukTMukTM/ukDNz4SO0MjL5QzM/g#top_anchor
type Client struct {
	client *lark.Lark
}

package feishu

import (
	"context"
	"fmt"

	"github.com/caarlos0/env/v6"
	"github.com/chyroc/lark"
)

func NewClient(appId, appSecret string) *Client {
	c := &Client{
		AppId:     appId,
		AppSecret: appSecret,
	}
	c.init()
	return c
}

func LoadClientFromEnv() (*Client, error) {
	c := &Client{}
	if err := env.Parse(c); err != nil {
		return nil, err
	}
	c.init()
	return c, nil
}

// 获取应用身份访问凭证 https://open.feishu.cn/document/ukTMukTMukTM/ukDNz4SO0MjL5QzM/g#top_anchor
type Client struct {
	AppId     string `env:"FEISHU_APP_ID"`
	AppSecret string `env:"FEISHU_APP_SECRET"`

	client *lark.Lark
}

func (c *Client) init() {
	c.client = lark.New(lark.WithAppCredential(c.AppId, c.AppSecret))
}

func (c *Client) Notify() {
	// us, resp, err := c.client.Contact.GetUser(context.Background(), &lark.GetUserReq{UserID: "ou_86791cfdf42c886435bdc65547d3ad8d"})
	// fmt.Println(*us.User, resp, err)
	mresp, resp, err := c.client.Message.Send().ToOpenID("ou_3551bad1f0b389385c6379e3df171e12").SendText(context.Background(), "测试个人通知")
	fmt.Println(mresp)
	fmt.Println(resp)
	fmt.Println(err)
}

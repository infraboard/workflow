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
	// 开发流程与工具介绍 https://open.feishu.cn/document/home/introduction-to-custom-app-development/self-built-application-development-process?lang=zh-CN
	AppId     string `env:"FEISHU_APP_ID"`
	AppSecret string `env:"FEISHU_APP_SECRET"`

	// 事件订阅概述 https://open.feishu.cn/document/ukTMukTMukTM/uUTNz4SN1MjL1UzM?lang=zh-CN
	// Encrypt Key 配置示例与解密: https://open.feishu.cn/document/ukTMukTMukTM/uYDNxYjL2QTM24iN0EjN/event-subscription-configure-/encrypt-key-encryption-configuration-case
	EncryptKey string `env:"FEISHU_ENCRYPT_KEY"`
	// 订阅事件时, 该自动会放置于token字段中, 用于校验发送方身份
	VerificationToken string `env:"FEISHU_VERIFICATION_TOKEN"`

	client *lark.Lark
}

func (c *Client) init() {
	c.client = lark.New(lark.WithAppCredential(c.AppId, c.AppSecret))
}

func (c *Client) Notify() {
	// us, resp, err := c.client.Contact.GetUser(context.Background(), &lark.GetUserReq{UserID: "ou_86791cfdf42c886435bdc65547d3ad8d"})
	// fmt.Println(*us.User, resp, err)
	mresp, resp, err := c.client.Message.Send().ToUserID("2fbc2b39").SendText(context.Background(), "测试个人通知")
	fmt.Println(mresp)
	fmt.Println(resp)
	fmt.Println(err)
}

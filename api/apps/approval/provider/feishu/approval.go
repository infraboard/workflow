// 对接开放平台审批API使用手册: https://open.feishu.cn/document/ukTMukTMukTM/uIjN4UjLyYDO14iM2gTN

package feishu

import (
	"context"

	"github.com/chyroc/lark"
)

func (c *Client) GetApproval(ctx context.Context, req *lark.GetApprovalReq) (*lark.GetApprovalResp, error) {
	resp, _, err := c.client.Approval.GetApproval(ctx, req)
	return resp, err
}

func (c *Client) CreateApprovalInstance(ctx context.Context, req *lark.CreateApprovalInstanceReq) (*lark.CreateApprovalInstanceResp, error) {
	resp, _, err := c.client.Approval.CreateApprovalInstance(ctx, req)
	return resp, err
}

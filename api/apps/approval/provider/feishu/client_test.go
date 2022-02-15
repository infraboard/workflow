package feishu_test

import (
	"context"
	"os"
	"testing"

	"github.com/chyroc/lark"
	"github.com/stretchr/testify/assert"

	"github.com/infraboard/workflow/api/apps/approval/provider/feishu"
)

func TestAuth(t *testing.T) {
	should := assert.New(t)
	client := feishu.NewClient(os.Getenv("FEISHU_APP_ID"), os.Getenv("FEISHU_APP_SECRET"))
	resp, err := client.GetApproval(context.Background(), &lark.GetApprovalReq{
		ApprovalCode: os.Getenv("ECS_APPROVAL_CODE"),
	})
	should.NoError(err)
	t.Log(resp)
}

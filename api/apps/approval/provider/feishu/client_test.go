package feishu_test

import (
	"testing"

	"github.com/infraboard/workflow/api/apps/approval/provider/feishu"
)

var (
	client *feishu.Client
)

func TestAuth(t *testing.T) {
	client.Notify()
}

func init() {
	c, err := feishu.LoadClientFromEnv()
	if err != nil {
		panic(err)
	}
	client = c
}

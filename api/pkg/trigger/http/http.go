package http

import (
	"errors"

	"github.com/infraboard/mcube/http/router"
	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/logger/zap"

	"github.com/infraboard/workflow/api/client"
	"github.com/infraboard/workflow/api/pkg"
	"github.com/infraboard/workflow/api/pkg/pipeline"
)

var (
	api = &handler{}
)

type handler struct {
	service pipeline.ServiceClient
	log     logger.Logger
}

// Registry 注册HTTP服务路由
func (h *handler) Registry(router router.SubRouter) {
	r := router.ResourceRouter("triggers")
	r.Auth(false)
	r.BasePath("triggers")
	r.Handle("POST", "/gitee", h.GiteeTrigger)
	r.Handle("POST", "/github", h.GitHubTrigger)
	r.Handle("POST", "/gitlab", h.GitLabTrigger)
}

func (h *handler) Config() error {
	client := client.C()
	if client == nil {
		return errors.New("grpc client not initial")
	}

	h.service = client.Pipeline()
	h.log = zap.L().Named("Trigger")
	return nil
}

func init() {
	pkg.RegistryHTTPV1("trigger", api)
}

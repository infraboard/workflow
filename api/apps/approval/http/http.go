package http

import (
	"github.com/infraboard/mcube/app"
	"github.com/infraboard/mcube/http/router"
	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/logger/zap"

	"github.com/infraboard/workflow/api/apps/approval"
)

var (
	api = &handler{}
)

type handler struct {
	service approval.ServiceServer

	log logger.Logger
}

// Registry 注册HTTP服务路由
func (h *handler) Registry(router router.SubRouter) {
	r := router.ResourceRouter("approval")
	r.Auth(true)
}

func (h *handler) Config() error {
	h.log = zap.L().Named(h.Name())
	h.service = app.GetGrpcApp(approval.AppName).(approval.ServiceServer)
	return nil
}

func (h *handler) Name() string {
	return approval.AppName
}

func init() {
	app.RegistryHttpApp(api)
}

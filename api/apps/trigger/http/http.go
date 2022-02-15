package http

import (
	"github.com/infraboard/mcube/app"
	"github.com/infraboard/mcube/http/router"
	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/logger/zap"

	"github.com/infraboard/workflow/api/apps/pipeline"
	"github.com/infraboard/workflow/api/apps/trigger"
)

var (
	api = &handler{}
)

type handler struct {
	service pipeline.ServiceServer
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
	h.service = app.GetGrpcApp(pipeline.AppName).(pipeline.ServiceServer)
	h.log = zap.L().Named("Trigger")
	return nil
}

func (h *handler) Name() string {
	return trigger.AppName
}

func init() {
	app.RegistryHttpApp(api)
}

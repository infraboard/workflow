package http

import (
	"github.com/infraboard/mcube/app"
	"github.com/infraboard/mcube/http/router"
	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/logger/zap"

	"github.com/infraboard/workflow/api/app/application"
)

var (
	api = &handler{}
)

type handler struct {
	service application.ServiceServer

	log logger.Logger
}

// Registry 注册HTTP服务路由
func (h *handler) Registry(router router.SubRouter) {
	r := router.ResourceRouter("application")
	r.Permission(true)
	r.BasePath("applications")
	r.Handle("POST", "/", h.CreateApplication)
	r.Handle("GET", "/", h.QueryApplication)
	r.Handle("GET", "/:id", h.DescribeApplication)
	r.Handle("PUT", "/:id", h.PutApplication)
	r.Handle("PATCH", "/:id", h.PatchApplication)
	r.Handle("DELETE", "/:name", h.DeleteApplication)

	r.BasePath("repo/projects")
	r.Handle("GET", "/", h.QuerySCMProject).DisableAuth()

	r.BasePath("triggers/scm/gitlab")
	r.Handle("POST", "/", h.GitLabHookHanler).DisableAuth()
}

func (h *handler) Config() error {
	h.log = zap.L().Named("Application")
	h.service = app.GetGrpcApp(application.AppName).(application.ServiceServer)
	return nil
}

func (h *handler) Name() string {
	return application.AppName
}

func init() {
	app.RegistryHttpApp(api)
}

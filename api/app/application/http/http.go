package http

import (
	"errors"

	"github.com/infraboard/mcube/http/router"
	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/logger/zap"

	"github.com/infraboard/workflow/api/app/application"
	"github.com/infraboard/workflow/api/client"
)

var (
	api = &handler{}
)

type handler struct {
	service application.ServiceClient

	log logger.Logger
}

// Registry 注册HTTP服务路由
func (h *handler) Registry(router router.SubRouter) {
	r := router.ResourceRouter("application")

	r.BasePath("applications")
	r.Handle("POST", "/", h.CreateApplication)
	r.Handle("GET", "/", h.QueryApplication)
	r.Handle("GET", "/:id", h.DescribeApplication)
	r.Handle("PUT", "/:id", h.PutApplication)
	r.Handle("PATCH", "/:id", h.PatchApplication)
	r.Handle("DELETE", "/:name", h.DeleteApplication)

	r.BasePath("repo/projects")
	r.Handle("GET", "/", h.QuerySCMProject)

	r.BasePath("triggers/scm/gitlab")
	r.Handle("POST", "/", h.GitLabHookHanler).DisableAuth()
}

func (h *handler) Config() error {
	h.log = zap.L().Named("Application")
	h.service = nil
	return nil
}

func (h *handler) Name() string {
	return application.AppName
}

func init() {
	app.RegistryHttpApp(api)
}

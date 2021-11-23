package http

import (
	"github.com/infraboard/mcube/http/label"
	"github.com/infraboard/mcube/http/router"

	"github.com/infraboard/mcube/app"
	"github.com/infraboard/workflow/api/app/deploy"
)

var (
	api = &handler{}
)

type handler struct {
	service deploy.ServiceServer
}

// Registry 注册HTTP服务路由
func (h *handler) Registry(router router.SubRouter) {
	r := router.ResourceRouter("deploy")

	r.BasePath("deploys")
	r.Handle("POST", "/", h.CreateApplicationDeploy).AddLabel(label.Create)
	r.Handle("GET", "/", h.QueryApplicationDeploy).AddLabel(label.List)
	r.Handle("GET", "/:id", h.DescribeApplicationDeploy).AddLabel(label.Get)
	r.Handle("DELETE", "/:id", h.DeleteApplicationDeploy).AddLabel(label.Delete)
}

func (h *handler) Config() error {
	h.service = nil
	return nil
}

func (h *handler) Name() string {
	return deploy.AppName
}

func init() {
	app.RegistryHttpApp(api)
}

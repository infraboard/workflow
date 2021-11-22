package http

import (
	"errors"

	"github.com/infraboard/mcube/http/router"
	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/logger/zap"

	"github.com/infraboard/workflow/api/app"
	"github.com/infraboard/workflow/api/app/template"
	"github.com/infraboard/workflow/api/client"
)

var (
	api = &handler{}
)

type handler struct {
	service template.ServiceClient
	log     logger.Logger
}

// Registry 注册HTTP服务路由
func (h *handler) Registry(router router.SubRouter) {
	r := router.ResourceRouter("templates")
	r.BasePath("templates")
	r.Handle("POST", "/", h.CreateTemplate)
	r.Handle("GET", "/", h.QueryTemplate)
	r.Handle("GET", "/:id", h.DescribeTemplate)
	r.Handle("PUT", "/:id", h.PutTemplate)
	r.Handle("PATCH", "/:id", h.PatchTemplate)
	r.Handle("DELETE", "/:id", h.DeleteTemplate)
}

func (h *handler) Config() error {
	h.service = nil

	h.log = zap.L().Named("Action")
	return nil
}

func init() {
	app.RegistryHttpApp(api)
}

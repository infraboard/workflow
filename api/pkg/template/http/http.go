package http

import (
	"errors"

	"github.com/infraboard/mcube/http/router"
	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/logger/zap"

	"github.com/infraboard/workflow/api/client"
	"github.com/infraboard/workflow/api/pkg"
	"github.com/infraboard/workflow/api/pkg/template"
)

var (
	api = &handler{log: zap.L().Named("Action")}
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
	r.Handle("DELETE", "/:id", h.DeleteTemplate)
}

func (h *handler) Config() error {
	client := client.C()
	if client == nil {
		return errors.New("grpc client not initial")
	}

	h.service = client.Template()
	return nil
}

func init() {
	pkg.RegistryHTTPV1("template", api)
}

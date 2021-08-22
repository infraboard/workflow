package http

import (
	"errors"

	"github.com/infraboard/mcube/http/router"

	"github.com/infraboard/workflow/api/client"
	"github.com/infraboard/workflow/api/pkg"
	"github.com/infraboard/workflow/api/pkg/application"
)

var (
	api = &handler{}
)

type handler struct {
	service application.ServiceClient
}

// Registry 注册HTTP服务路由
func (h *handler) Registry(router router.SubRouter) {
	r := router.ResourceRouter("application")

	r.BasePath("applications")
	r.Handle("POST", "/", h.CreateApplication)
	r.Handle("GET", "/", h.QueryApplication)

	r.BasePath("repo/projects")
	r.Handle("GET", "/", h.QueryRepoProject)
}

func (h *handler) Config() error {
	client := client.C()
	if client == nil {
		return errors.New("grpc client not initial")
	}

	h.service = client.Application()
	return nil
}

func init() {
	pkg.RegistryHTTPV1("application", api)
}

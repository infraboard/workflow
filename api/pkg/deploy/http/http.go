package http

import (
	"errors"

	"github.com/infraboard/mcube/http/label"
	"github.com/infraboard/mcube/http/router"

	"github.com/infraboard/workflow/api/client"
	"github.com/infraboard/workflow/api/pkg"
	"github.com/infraboard/workflow/api/pkg/deploy"
)

var (
	api = &handler{}
)

type handler struct {
	service deploy.ServiceClient
}

// Registry 注册HTTP服务路由
func (h *handler) Registry(router router.SubRouter) {
	r := router.ResourceRouter("deploy")

	r.BasePath("deploys")
	r.Handle("POST", "/", h.CreateApplicationDeploy).AddLabel(label.Create)
	r.Handle("GET", "/", h.QueryApplicationDeploy).AddLabel(label.List)
}

func (h *handler) Config() error {
	client := client.C()
	if client == nil {
		return errors.New("grpc client not initial")
	}

	h.service = client.Deploy()
	return nil
}

func init() {
	pkg.RegistryHTTPV1("deploy", api)
}

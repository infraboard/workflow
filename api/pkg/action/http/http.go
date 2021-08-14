package http

import (
	"errors"

	"github.com/infraboard/mcube/http/label"
	"github.com/infraboard/mcube/http/router"
	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/logger/zap"

	"github.com/infraboard/workflow/api/client"
	"github.com/infraboard/workflow/api/pkg"
	"github.com/infraboard/workflow/api/pkg/action"
)

var (
	api = &handler{log: zap.L().Named("Action")}
)

type handler struct {
	service action.ServiceClient
	log     logger.Logger
}

// Registry 注册HTTP服务路由
func (h *handler) Registry(router router.SubRouter) {
	r := router.ResourceRouter("actions")
	r.BasePath("actions")
	r.Handle("POST", "/", h.CreateAction).AddLabel(label.Create)
	r.Handle("GET", "/", h.QueryAction).AddLabel(label.List)
	r.Handle("GET", "/:key", h.DescribeAction).AddLabel(label.Get)
	r.Handle("DELETE", "/:key", h.DeleteNamespaceAction).AddLabel(label.Delete)
	r.BasePath("global_actions")
	r.Handle("POST", "/", h.CreateGlobalAction).AddLabel(label.Create)
	r.Handle("DELETE", "/:name", h.DeleteGlobalAction).AddLabel(label.Delete)
	r.BasePath("runners")
	r.Handle("GET", "/", h.QueryRunner).AddLabel(label.List)
}

func (h *handler) Config() error {
	client := client.C()
	if client == nil {
		return errors.New("grpc client not initial")
	}

	h.service = client.Action()
	return nil
}

func init() {
	pkg.RegistryHTTPV1("action", api)
}

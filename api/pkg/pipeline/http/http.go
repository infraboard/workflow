package http

import (
	"errors"

	"github.com/infraboard/mcube/http/label"
	"github.com/infraboard/mcube/http/router"

	"github.com/infraboard/workflow/api/client"
	"github.com/infraboard/workflow/api/pkg"
	"github.com/infraboard/workflow/api/pkg/pipeline"
)

var (
	api = &handler{}
)

type handler struct {
	service pipeline.ServiceClient
}

// Registry 注册HTTP服务路由
func (h *handler) Registry(router router.SubRouter) {
	r := router.ResourceRouter("pipeline")
	r.BasePath("pipelines")
	r.Handle("POST", "/", h.CreatePipeline).AddLabel(label.Create)
	r.Handle("GET", "/", h.QueryPipeline).AddLabel(label.List)
	r.Handle("GET", "/:id", h.DescribePipeline).AddLabel(label.Get)
	r.Handle("DELETE", "/:id", h.DeletePipeline).AddLabel(label.Delete)

	r.BasePath("actions")
	r.Handle("POST", "/", h.CreateAction).AddLabel(label.Create)
	r.Handle("GET", "/", h.QueryAction).AddLabel(label.List)
	r.Handle("DELETE", "/:name", h.DeleteAction).AddLabel(label.Delete)

}

func (h *handler) Config() error {
	client := client.C()
	if client == nil {
		return errors.New("grpc client not initial")
	}

	h.service = client.Pipeline()
	return nil
}

func init() {
	pkg.RegistryHTTPV1("pipeline", api)
}

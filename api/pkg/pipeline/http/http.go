package http

import (
	"errors"

	"github.com/infraboard/mcube/http/label"
	"github.com/infraboard/mcube/http/router"
	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/logger/zap"

	"github.com/infraboard/workflow/api/client"
	"github.com/infraboard/workflow/api/pkg"
	"github.com/infraboard/workflow/api/pkg/pipeline"
)

var (
	API = &handler{log: zap.L().Named("Pipeline")}
)

type handler struct {
	service pipeline.ServiceClient
	log     logger.Logger
	proxy   *Proxy
}

// Registry 注册HTTP服务路由
func (h *handler) Registry(router router.SubRouter) {
	r := router.ResourceRouter("pipeline")
	r.BasePath("pipelines")
	r.Handle("POST", "/", h.CreatePipeline).AddLabel(label.Create)
	r.Handle("GET", "/", h.QueryPipeline).AddLabel(label.List)
	r.Handle("GET", "/:id", h.DescribePipeline).AddLabel(label.Get)
	r.Handle("DELETE", "/:id", h.DeletePipeline).AddLabel(label.Delete)
	r.Handle("GET", "/:id/watch_check", h.WatchPipelineCheck).AddLabel(label.Get)
	r.BasePath("websocket")
	r.Handle("GET", "pipelines/:id/watch", h.WatchPipeline).AddLabel(label.Get)

	r.BasePath("steps")
	r.Handle("GET", "/", h.QueryStep).AddLabel(label.List)
	r.Handle("POST", "/", h.CreateStep).AddLabel(label.Create)
	r.Handle("GET", "/:id", h.DescribeStep).AddLabel(label.Get)
	r.Handle("DELETE", "/:id", h.DeleteStep).AddLabel(label.Delete)
	r.Handle("POST", "/:id/audit", h.AuditStep).AddLabel(label.Update)
	r.Handle("POST", "/:id/cancel", h.CancelStep).AddLabel(label.Update)

	r.BasePath("actions")
	r.Handle("POST", "/", h.CreateAction).AddLabel(label.Create)
	r.Handle("GET", "/", h.QueryAction).AddLabel(label.List)
	r.Handle("GET", "/:name", h.DescribeAction).AddLabel(label.Get)
	r.Handle("DELETE", "/:name", h.DeleteNamespaceAction).AddLabel(label.Delete)
	r.BasePath("global_actions")
	r.Handle("POST", "/", h.CreateGlobalAction).AddLabel(label.Create)
	r.Handle("DELETE", "/:name", h.DeleteGlobalAction).AddLabel(label.Delete)
}

func (h *handler) Config() error {
	client := client.C()
	if client == nil {
		return errors.New("grpc client not initial")
	}

	h.service = client.Pipeline()
	h.proxy = NewProxy()
	return nil
}

func init() {
	pkg.RegistryHTTPV1("pipeline", API)
}

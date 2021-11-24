package http

import (
	"github.com/infraboard/mcube/app"
	"github.com/infraboard/mcube/http/label"
	"github.com/infraboard/mcube/http/router"
	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/logger/zap"

	"github.com/infraboard/workflow/api/app/pipeline"
)

var (
	api = &handler{}
)

type handler struct {
	service pipeline.ServiceServer
	log     logger.Logger
	proxy   *Proxy
}

// Registry 注册HTTP服务路由
func (h *handler) Registry(router router.SubRouter) {
	r := router.ResourceRouter("pipeline")
	r.Permission(true)
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

	r.BasePath("variable_templates")
	r.Handle("GET", "/", h.QueryVariableTemplate).AddLabel(label.List)
	r.BasePath("enums")
	r.Handle("GET", "/step_status", h.QueryStepStatusEnum).AddLabel(label.List)
}

func (h *handler) Config() error {
	h.service = nil
	h.proxy = NewProxy()

	h.log = zap.L().Named("Pipeline")
	return nil
}

func (h *handler) Name() string {
	return pipeline.AppName
}

func init() {
	app.RegistryHttpApp(api)
}

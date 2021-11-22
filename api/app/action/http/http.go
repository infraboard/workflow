package http

import (
	"errors"

	"github.com/infraboard/mcube/http/label"
	"github.com/infraboard/mcube/http/router"
	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/logger/zap"

	"github.com/infraboard/workflow/api/app/action"
	"github.com/infraboard/workflow/api/client"
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
	r.Handle("PUT", "/:key", h.UpdateAction).AddLabel()
	r.Handle("DELETE", "/:key", h.DeleteAction).AddLabel(label.Delete)

	r.BasePath("runners")
	r.Handle("GET", "/", h.QueryRunner).AddLabel(label.List)
}

func (h *handler) Config() error {
	h.service = app.GetGrpcApp(action.AppName).(action.ServiceServer)
	return nil
}

func (h *handler) Name() string {
	return action.AppName
}

func init() {
	app.RegistryHttpApp(api)
}

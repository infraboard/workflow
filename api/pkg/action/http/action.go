package http

import (
	"fmt"
	"net/http"

	"github.com/infraboard/keyauth/pkg/token"
	"github.com/infraboard/mcube/grpc/gcontext"
	"github.com/infraboard/mcube/http/context"
	"github.com/infraboard/mcube/http/request"
	"github.com/infraboard/mcube/http/response"
	"github.com/infraboard/mcube/pb/resource"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/infraboard/workflow/api/pkg/action"
	"github.com/infraboard/workflow/node/controller/step/runner/docker"
	"github.com/infraboard/workflow/node/controller/step/runner/k8s"
	"github.com/infraboard/workflow/node/controller/step/runner/local"
)

// Action
func (h *handler) CreateAction(w http.ResponseWriter, r *http.Request) {
	ctx, err := gcontext.NewGrpcOutCtxFromHTTPRequest(r)
	if err != nil {
		response.Failed(w, err)
		return
	}

	hc := context.GetContext(r)
	tk, ok := hc.AuthInfo.(*token.Token)
	if !ok {
		response.Failed(w, fmt.Errorf("auth info is not an *token.Token"))
		return
	}

	req := action.NewCreateActionRequest()
	if err := request.GetDataFromRequest(r, req); err != nil {
		response.Failed(w, err)
		return
	}
	req.VisiableMode = resource.VisiableMode_NAMESPACE
	req.UpdateOwner(tk)

	var header, trailer metadata.MD
	ins, err := h.service.CreateAction(
		ctx.Context(),
		req,
		grpc.Header(&header),
		grpc.Trailer(&trailer),
	)
	if err != nil {
		response.Failed(w, gcontext.NewExceptionFromTrailer(trailer, err))
		return
	}

	ins.InitNil()
	response.Success(w, ins)
}

func (h *handler) UpdateAction(w http.ResponseWriter, r *http.Request) {
	ctx, err := gcontext.NewGrpcOutCtxFromHTTPRequest(r)
	if err != nil {
		response.Failed(w, err)
		return
	}

	hc := context.GetContext(r)
	name, version := action.ParseActionKey(hc.PS.ByName("key"))
	req := action.NewUpdateActionRequest()
	if err := request.GetDataFromRequest(r, req); err != nil {
		response.Failed(w, err)
		return
	}
	req.Name = name
	req.Version = version

	var header, trailer metadata.MD
	ins, err := h.service.UpdateAction(
		ctx.Context(),
		req,
		grpc.Header(&header),
		grpc.Trailer(&trailer),
	)
	if err != nil {
		response.Failed(w, gcontext.NewExceptionFromTrailer(trailer, err))
		return
	}

	ins.InitNil()
	response.Success(w, ins)
}

func (h *handler) QueryAction(w http.ResponseWriter, r *http.Request) {
	ctx, err := gcontext.NewGrpcOutCtxFromHTTPRequest(r)
	if err != nil {
		response.Failed(w, err)
		return
	}

	hc := context.GetContext(r)
	tk, ok := hc.AuthInfo.(*token.Token)
	if !ok {
		response.Failed(w, fmt.Errorf("auth info is not an *token.Token"))
		return
	}

	page := request.NewPageRequestFromHTTP(r)
	req := action.NewQueryActionRequest(page)
	req.Namespace = tk.Namespace

	var header, trailer metadata.MD
	actions, err := h.service.QueryAction(
		ctx.Context(),
		req,
		grpc.Header(&header),
		grpc.Trailer(&trailer),
	)
	if err != nil {
		response.Failed(w, gcontext.NewExceptionFromTrailer(trailer, err))
		return
	}

	// 避免前端处理null
	actions.InitNil()
	response.Success(w, actions)
}

func (h *handler) DescribeAction(w http.ResponseWriter, r *http.Request) {
	ctx, err := gcontext.NewGrpcOutCtxFromHTTPRequest(r)
	if err != nil {
		response.Failed(w, err)
		return
	}

	hc := context.GetContext(r)
	name, version := action.ParseActionKey(hc.PS.ByName("key"))
	req := action.NewDescribeActionRequest(name, version)

	var header, trailer metadata.MD
	ins, err := h.service.DescribeAction(
		ctx.Context(),
		req,
		grpc.Header(&header),
		grpc.Trailer(&trailer),
	)
	if err != nil {
		response.Failed(w, gcontext.NewExceptionFromTrailer(trailer, err))
		return
	}

	ins.InitNil()
	response.Success(w, ins)
}

func (h *handler) DeleteAction(w http.ResponseWriter, r *http.Request) {
	ctx, err := gcontext.NewGrpcOutCtxFromHTTPRequest(r)
	if err != nil {
		response.Failed(w, err)
		return
	}

	hc := context.GetContext(r)
	name, version := action.ParseActionKey(hc.PS.ByName("key"))
	req := action.NewDeleteActionRequest(name, version)

	tk, ok := hc.AuthInfo.(*token.Token)
	if !ok {
		response.Failed(w, fmt.Errorf("auth info is not an *token.Token"))
		return
	}
	req.Namespace = tk.Namespace

	var header, trailer metadata.MD
	action, err := h.service.DeleteAction(
		ctx.Context(),
		req,
		grpc.Header(&header),
		grpc.Trailer(&trailer),
	)
	if err != nil {
		response.Failed(w, gcontext.NewExceptionFromTrailer(trailer, err))
		return
	}
	response.Success(w, action)
}

func (h *handler) QueryRunner(w http.ResponseWriter, r *http.Request) {
	ins := NewRunnerParamDescSet()
	ins.Add(action.RUNNER_TYPE_DOCKER, docker.ParamsDesc())
	ins.Add(action.RUNNER_TYPE_K8s, k8s.ParamsDesc())
	ins.Add(action.RUNNER_TYPE_LOCAL, local.ParamsDesc())
	response.Success(w, ins)
}

func NewRunnerParamDescSet() *RunnerParamDescSet {
	return &RunnerParamDescSet{
		Items: []*RunnerParamDesc{},
	}
}

type RunnerParamDescSet struct {
	Items []*RunnerParamDesc `json:"items"`
}

func (s *RunnerParamDescSet) Add(t action.RUNNER_TYPE, desc []*action.RunParamDesc) {
	s.Items = append(s.Items, &RunnerParamDesc{
		Type:      t,
		ParamDesc: desc,
	})
}

type RunnerParamDesc struct {
	Type      action.RUNNER_TYPE     `json:"type"`
	ParamDesc []*action.RunParamDesc `json:"param_desc"`
}

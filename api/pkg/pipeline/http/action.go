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

	"github.com/infraboard/workflow/api/pkg/pipeline"
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

	req := pipeline.NewCreateActionRequest()
	if err := request.GetDataFromRequest(r, req); err != nil {
		response.Failed(w, err)
		return
	}
	req.VisiableMode = resource.VisiableMode_NAMESPACE

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
	response.Success(w, ins)
}

func (h *handler) QueryAction(w http.ResponseWriter, r *http.Request) {
	ctx, err := gcontext.NewGrpcOutCtxFromHTTPRequest(r)
	if err != nil {
		response.Failed(w, err)
		return
	}

	page := request.NewPageRequestFromHTTP(r)
	req := pipeline.NewQueryActionRequest(page)

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
	response.Success(w, actions)
}

func (h *handler) DescribeAction(w http.ResponseWriter, r *http.Request) {
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

	fmt.Println(hc.PS.ByName("key"))
	name, version := pipeline.ParseActionKey(hc.PS.ByName("key"))
	req := pipeline.NewDescribeActionRequest(tk.Namespace, name, version)

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
	response.Success(w, ins)
}

func (h *handler) DeleteNamespaceAction(w http.ResponseWriter, r *http.Request) {
	ctx, err := gcontext.NewGrpcOutCtxFromHTTPRequest(r)
	if err != nil {
		response.Failed(w, err)
		return
	}

	hc := context.GetContext(r)
	name, version := pipeline.ParseActionKey(hc.PS.ByName("key"))
	req := pipeline.NewDeleteActionRequest(name, version)
	req.VisiableMode = resource.VisiableMode_NAMESPACE

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

func (h *handler) DeleteGlobalAction(w http.ResponseWriter, r *http.Request) {
	ctx, err := gcontext.NewGrpcOutCtxFromHTTPRequest(r)
	if err != nil {
		response.Failed(w, err)
		return
	}

	hc := context.GetContext(r)
	name, version := pipeline.ParseActionKey(hc.PS.ByName("key"))
	req := pipeline.NewDeleteActionRequest(name, version)
	req.VisiableMode = resource.VisiableMode_GLOBAL

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

// Action
func (h *handler) CreateGlobalAction(w http.ResponseWriter, r *http.Request) {
	ctx, err := gcontext.NewGrpcOutCtxFromHTTPRequest(r)
	if err != nil {
		response.Failed(w, err)
		return
	}

	req := pipeline.NewCreateActionRequest()
	if err := request.GetDataFromRequest(r, req); err != nil {
		response.Failed(w, err)
		return
	}
	req.VisiableMode = resource.VisiableMode_GLOBAL

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
	response.Success(w, ins)
}

func (h *handler) QueryRunner(w http.ResponseWriter, r *http.Request) {
	ins := NewRunnerParamDescSet()
	ins.Add(pipeline.RUNNER_TYPE_DOCKER, docker.ParamsDesc())
	ins.Add(pipeline.RUNNER_TYPE_K8s, k8s.ParamsDesc())
	ins.Add(pipeline.RUNNER_TYPE_LOCAL, local.ParamsDesc())
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

func (s *RunnerParamDescSet) Add(t pipeline.RUNNER_TYPE, desc []*pipeline.RunParamDesc) {
	s.Items = append(s.Items, &RunnerParamDesc{
		Type:      t,
		ParamDesc: desc,
	})
}

type RunnerParamDesc struct {
	Type      pipeline.RUNNER_TYPE     `json:"type"`
	ParamDesc []*pipeline.RunParamDesc `json:"param_desc"`
}

package http

import (
	"fmt"
	"net/http"

	"github.com/infraboard/keyauth/app/token"
	"github.com/infraboard/mcube/grpc/gcontext"
	"github.com/infraboard/mcube/http/context"
	"github.com/infraboard/mcube/http/request"
	"github.com/infraboard/mcube/http/response"

	"github.com/infraboard/workflow/api/app/pipeline"
)

// Action
func (h *handler) CreateStep(w http.ResponseWriter, r *http.Request) {
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

	req := pipeline.NewCreateStepRequest()
	if err := request.GetDataFromRequest(r, req); err != nil {
		response.Failed(w, err)
		return
	}
	req.Namespace = tk.Namespace

	ins, err := h.service.CreateStep(
		ctx.Context(),
		req,
	)
	if err != nil {
		response.Failed(w, err)
		return
	}
	response.Success(w, ins)
}

func (h *handler) QueryStep(w http.ResponseWriter, r *http.Request) {
	ctx, err := gcontext.NewGrpcOutCtxFromHTTPRequest(r)
	if err != nil {
		response.Failed(w, err)
		return
	}

	page := request.NewPageRequestFromHTTP(r)
	req := pipeline.NewQueryStepRequest()
	req.Page = &page.PageRequest

	dommains, err := h.service.QueryStep(
		ctx.Context(),
		req,
	)
	if err != nil {
		response.Failed(w, err)
		return
	}
	response.Success(w, dommains)
}

func (h *handler) DescribeStep(w http.ResponseWriter, r *http.Request) {
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

	req := pipeline.NewDescribeStepRequestWithKey(hc.PS.ByName("id"))
	req.Namespace = tk.Namespace

	dommains, err := h.service.DescribeStep(
		ctx.Context(),
		req,
	)
	if err != nil {
		response.Failed(w, err)
		return
	}
	response.Success(w, dommains)
}

func (h *handler) DeleteStep(w http.ResponseWriter, r *http.Request) {
	ctx, err := gcontext.NewGrpcOutCtxFromHTTPRequest(r)
	if err != nil {
		response.Failed(w, err)
		return
	}

	hc := context.GetContext(r)
	req := pipeline.NewDeleteStepRequestWithKey(hc.PS.ByName("id"))

	dommains, err := h.service.DeleteStep(
		ctx.Context(),
		req,
	)
	if err != nil {
		response.Failed(w, err)
		return
	}
	response.Success(w, dommains)
}

func (h *handler) AuditStep(w http.ResponseWriter, r *http.Request) {
	ctx, err := gcontext.NewGrpcOutCtxFromHTTPRequest(r)
	if err != nil {
		response.Failed(w, err)
		return
	}

	hc := context.GetContext(r)
	req := pipeline.NewAuditStepRequest()
	if err := request.GetDataFromRequest(r, req); err != nil {
		response.Failed(w, err)
		return
	}
	req.Key = hc.PS.ByName("id")

	dommains, err := h.service.AuditStep(
		ctx.Context(),
		req,
	)
	if err != nil {
		response.Failed(w, err)
		return
	}
	response.Success(w, dommains)
}

func (h *handler) CancelStep(w http.ResponseWriter, r *http.Request) {
	ctx, err := gcontext.NewGrpcOutCtxFromHTTPRequest(r)
	if err != nil {
		response.Failed(w, err)
		return
	}

	hc := context.GetContext(r)
	req := pipeline.NewCancelStepRequestWithKey(hc.PS.ByName("id"))

	dommains, err := h.service.CancelStep(
		ctx.Context(),
		req,
	)
	if err != nil {
		response.Failed(w, err)
		return
	}
	response.Success(w, dommains)
}

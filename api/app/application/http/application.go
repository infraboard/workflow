package http

import (
	"fmt"
	"net/http"

	"github.com/infraboard/keyauth/app/token"
	"github.com/infraboard/mcube/http/context"
	"github.com/infraboard/mcube/http/request"
	"github.com/infraboard/mcube/http/response"
	pb "github.com/infraboard/mcube/pb/request"

	"github.com/infraboard/workflow/api/app/application"
)

func (h *handler) CreateApplication(w http.ResponseWriter, r *http.Request) {
	ctx := context.GetContext(r)
	tk := ctx.AuthInfo.(*token.Token)

	req := application.NewCreateApplicationRequest()
	if err := request.GetDataFromRequest(r, req); err != nil {
		response.Failed(w, err)
		return
	}
	req.UpdateOwner(tk)

	ins, err := h.service.CreateApplication(
		r.Context(),
		req,
	)
	if err != nil {
		response.Failed(w, err)
		return
	}
	response.Success(w, ins)
}

func (h *handler) QueryApplication(w http.ResponseWriter, r *http.Request) {
	ctx := context.GetContext(r)
	tk := ctx.AuthInfo.(*token.Token)

	page := request.NewPageRequestFromHTTP(r)
	req := application.NewQueryApplicationRequest(page)
	req.Domain = tk.Domain
	req.Namespace = tk.Namespace

	dommains, err := h.service.QueryApplication(
		r.Context(),
		req,
	)
	if err != nil {
		response.Failed(w, err)
		return
	}
	response.Success(w, dommains)
}

func (h *handler) DescribeApplication(w http.ResponseWriter, r *http.Request) {
	ctx := context.GetContext(r)
	req := application.NewDescribeApplicationRequestWithID(ctx.PS.ByName("id"))
	ins, err := h.service.DescribeApplication(
		r.Context(),
		req,
	)
	if err != nil {
		response.Failed(w, err)
		return
	}
	ins.Desensitize()
	response.Success(w, ins)
}

func (h *handler) DeleteApplication(w http.ResponseWriter, r *http.Request) {
	ctx := context.GetContext(r)
	tk, ok := ctx.AuthInfo.(*token.Token)
	if !ok {
		response.Failed(w, fmt.Errorf("auth info is not an *token.Token"))
		return
	}

	req := application.NewDeleteApplicationRequest(tk.Namespace, ctx.PS.ByName("name"))

	action, err := h.service.DeleteApplication(
		r.Context(),
		req,
	)
	if err != nil {
		response.Failed(w, err)
		return
	}
	response.Success(w, action)
}

func (h *handler) PutApplication(w http.ResponseWriter, r *http.Request) {
	ctx := context.GetContext(r)
	tk, ok := ctx.AuthInfo.(*token.Token)
	if !ok {
		response.Failed(w, fmt.Errorf("auth info is not an *token.Token"))
		return
	}

	req := application.NewUpdateApplicationRequest(ctx.PS.ByName("id"))
	req.UpdateBy = tk.Account
	if err := request.GetDataFromRequest(r, req.Data); err != nil {
		response.Failed(w, err)
		return
	}

	ins, err := h.service.UpdateApplication(
		r.Context(),
		req,
	)
	if err != nil {
		response.Failed(w, err)
		return
	}
	ins.Desensitize()

	response.Success(w, ins)
	return
}

func (h *handler) PatchApplication(w http.ResponseWriter, r *http.Request) {
	ctx := context.GetContext(r)
	tk, ok := ctx.AuthInfo.(*token.Token)
	if !ok {
		response.Failed(w, fmt.Errorf("auth info is not an *token.Token"))
		return
	}

	req := application.NewUpdateApplicationRequest(ctx.PS.ByName("id"))
	req.UpdateMode = pb.UpdateMode_PATCH
	req.UpdateBy = tk.Account
	if err := request.GetDataFromRequest(r, req.Data); err != nil {
		response.Failed(w, err)
		return
	}

	ins, err := h.service.UpdateApplication(
		r.Context(),
		req,
	)
	if err != nil {
		response.Failed(w, err)
		return
	}
	ins.Desensitize()

	response.Success(w, ins)
	return
}

package http

import (
	"net/http"

	"github.com/infraboard/keyauth/app/token"
	"github.com/infraboard/mcube/http/context"
	"github.com/infraboard/mcube/http/request"
	"github.com/infraboard/mcube/http/response"
	pb "github.com/infraboard/mcube/pb/request"

	"github.com/infraboard/workflow/api/apps/template"
)

// Action
func (h *handler) CreateTemplate(w http.ResponseWriter, r *http.Request) {
	ctx := context.GetContext(r)
	tk := ctx.AuthInfo.(*token.Token)

	req := template.NewCreateTemplateRequest()
	if err := request.GetDataFromRequest(r, req); err != nil {
		response.Failed(w, err)
		return
	}
	req.UpdateOwner(tk)

	ins, err := h.service.CreateTemplate(
		r.Context(),
		req,
	)
	if err != nil {
		response.Failed(w, err)
		return
	}
	response.Success(w, ins)
}

func (h *handler) QueryTemplate(w http.ResponseWriter, r *http.Request) {
	ctx := context.GetContext(r)
	tk := ctx.AuthInfo.(*token.Token)

	page := request.NewPageRequestFromHTTP(r)
	req := template.NewQueryTemplateRequest(page)
	req.Namespace = tk.Namespace

	actions, err := h.service.QueryTemplate(
		r.Context(),
		req,
	)
	if err != nil {
		response.Failed(w, err)
		return
	}
	response.Success(w, actions)
}

func (h *handler) DescribeTemplate(w http.ResponseWriter, r *http.Request) {
	ctx := context.GetContext(r)
	req := template.NewDescribeTemplateRequestWithID(ctx.PS.ByName("id"))

	ins, err := h.service.DescribeTemplate(
		r.Context(),
		req,
	)
	if err != nil {
		response.Failed(w, err)
		return
	}
	response.Success(w, ins)
}

func (h *handler) DeleteTemplate(w http.ResponseWriter, r *http.Request) {
	ctx := context.GetContext(r)
	req := template.NewDeleteTemplateRequestWithID(ctx.PS.ByName("id"))

	action, err := h.service.DeleteTemplate(
		r.Context(),
		req,
	)
	if err != nil {
		response.Failed(w, err)
		return
	}
	response.Success(w, action)
}

func (h *handler) PutTemplate(w http.ResponseWriter, r *http.Request) {
	ctx := context.GetContext(r)
	tk := ctx.AuthInfo.(*token.Token)

	req := template.NewUpdateTemplateRequest(ctx.PS.ByName("id"))
	req.UpdateBy = tk.Account
	if err := request.GetDataFromRequest(r, req.Data); err != nil {
		response.Failed(w, err)
		return
	}

	ins, err := h.service.UpdateTemplate(
		r.Context(),
		req,
	)
	if err != nil {
		response.Failed(w, err)
		return
	}

	response.Success(w, ins)
}

func (h *handler) PatchTemplate(w http.ResponseWriter, r *http.Request) {
	ctx := context.GetContext(r)
	tk := ctx.AuthInfo.(*token.Token)

	req := template.NewUpdateTemplateRequest(ctx.PS.ByName("id"))
	req.UpdateMode = pb.UpdateMode_PATCH
	req.UpdateBy = tk.Account
	if err := request.GetDataFromRequest(r, req.Data); err != nil {
		response.Failed(w, err)
		return
	}

	ins, err := h.service.UpdateTemplate(
		r.Context(),
		req,
	)
	if err != nil {
		response.Failed(w, err)
		return
	}

	response.Success(w, ins)
}

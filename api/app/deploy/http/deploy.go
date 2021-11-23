package http

import (
	"net/http"

	"github.com/infraboard/keyauth/app/token"
	"github.com/infraboard/mcube/http/context"
	"github.com/infraboard/mcube/http/request"
	"github.com/infraboard/mcube/http/response"

	"github.com/infraboard/workflow/api/app/deploy"
)

func (h *handler) CreateApplicationDeploy(w http.ResponseWriter, r *http.Request) {
	ctx := context.GetContext(r)
	tk := ctx.AuthInfo.(*token.Token)

	req := deploy.NewCreateApplicationDeployRequest()
	if err := request.GetDataFromRequest(r, req); err != nil {
		response.Failed(w, err)
		return
	}
	req.UpdateOwner(tk)

	ins, err := h.service.CreateApplicationDeploy(
		r.Context(),
		req,
	)
	if err != nil {
		response.Failed(w, err)
		return
	}
	response.Success(w, ins)
}

func (h *handler) QueryApplicationDeploy(w http.ResponseWriter, r *http.Request) {
	ctx := context.GetContext(r)
	tk := ctx.AuthInfo.(*token.Token)

	page := request.NewPageRequestFromHTTP(r)
	req := deploy.NewQueryApplicationDeployRequest(page)
	req.Domain = tk.Domain
	req.Namespace = tk.Namespace

	dommains, err := h.service.QueryApplicationDeploy(
		r.Context(),
		req,
	)
	if err != nil {
		response.Failed(w, err)
		return
	}
	response.Success(w, dommains)
}

func (h *handler) DescribeApplicationDeploy(w http.ResponseWriter, r *http.Request) {
	ctx := context.GetContext(r)
	req := deploy.NewDescribeApplicationDeployRequestWithID(ctx.PS.ByName("id"))
	ins, err := h.service.DescribeApplicationDeploy(
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

func (h *handler) DeleteApplicationDeploy(w http.ResponseWriter, r *http.Request) {
	ctx := context.GetContext(r)
	tk := ctx.AuthInfo.(*token.Token)

	req := deploy.NewDeleteApplicationDeployRequest(tk.Namespace, ctx.PS.ByName("id"))

	action, err := h.service.DeleteApplicationDeploy(
		r.Context(),
		req,
	)
	if err != nil {
		response.Failed(w, err)
		return
	}
	response.Success(w, action)
}

package http

import (
	"fmt"
	"io"
	"net/http"

	"github.com/infraboard/keyauth/pkg/token"
	"github.com/infraboard/mcube/grpc/gcontext"
	"github.com/infraboard/mcube/http/context"
	"github.com/infraboard/mcube/http/request"
	"github.com/infraboard/mcube/http/response"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/infraboard/workflow/api/pkg/pipeline"
)

func (h *handler) CreatePipeline(w http.ResponseWriter, r *http.Request) {
	ctx, err := gcontext.NewGrpcOutCtxFromHTTPRequest(r)
	if err != nil {
		response.Failed(w, err)
		return
	}

	req := pipeline.NewCreatePipelineRequest()
	if err := request.GetDataFromRequest(r, req); err != nil {
		response.Failed(w, err)
		return
	}

	var header, trailer metadata.MD
	ins, err := h.service.CreatePipeline(
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

func (h *handler) QueryPipeline(w http.ResponseWriter, r *http.Request) {
	ctx, err := gcontext.NewGrpcOutCtxFromHTTPRequest(r)
	if err != nil {
		response.Failed(w, err)
		return
	}

	page := request.NewPageRequestFromHTTP(r)
	req := pipeline.NewQueryPipelineRequest()
	req.Page = &page.PageRequest

	var header, trailer metadata.MD
	dommains, err := h.service.QueryPipeline(
		ctx.Context(),
		req,
		grpc.Header(&header),
		grpc.Trailer(&trailer),
	)
	if err != nil {
		response.Failed(w, gcontext.NewExceptionFromTrailer(trailer, err))
		return
	}
	response.Success(w, dommains)
}

func (h *handler) DescribePipeline(w http.ResponseWriter, r *http.Request) {
	ctx, err := gcontext.NewGrpcOutCtxFromHTTPRequest(r)
	if err != nil {
		response.Failed(w, err)
		return
	}

	hc := context.GetContext(r)
	req := pipeline.NewDescribePipelineRequestWithID(hc.PS.ByName("id"))

	var header, trailer metadata.MD
	dommains, err := h.service.DescribePipeline(
		ctx.Context(),
		req,
		grpc.Header(&header),
		grpc.Trailer(&trailer),
	)
	if err != nil {
		response.Failed(w, gcontext.NewExceptionFromTrailer(trailer, err))
		return
	}
	response.Success(w, dommains)
}

// pipeline删除时,除了删除pipeline对象本身而外，还需要删除该pipeline下的所有step
func (h *handler) DeletePipeline(w http.ResponseWriter, r *http.Request) {
	ctx, err := gcontext.NewGrpcOutCtxFromHTTPRequest(r)
	if err != nil {
		response.Failed(w, err)
		return
	}

	hc := context.GetContext(r)
	req := pipeline.NewDeletePipelineRequestWithID(hc.PS.ByName("id"))

	var header, trailer metadata.MD
	dommains, err := h.service.DeletePipeline(
		ctx.Context(),
		req,
		grpc.Header(&header),
		grpc.Trailer(&trailer),
	)
	if err != nil {
		response.Failed(w, gcontext.NewExceptionFromTrailer(trailer, err))
		return
	}
	response.Success(w, dommains)
}

func (h *handler) WatchPipelineCheck(w http.ResponseWriter, r *http.Request) {
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

	req := pipeline.NewWatchPipelineRequestByID("", hc.PS.ByName("id"))
	req.Namespace = tk.Namespace
	req.DryRun = true

	var header, trailer metadata.MD
	stream, err := h.service.WatchPipeline(
		ctx.Context(),
		req,
		grpc.Header(&header),
		grpc.Trailer(&trailer),
	)

	if err != nil {
		response.Failed(w, gcontext.NewExceptionFromTrailer(trailer, err))
		return
	}

	_, err = stream.Recv()
	if err != nil && err != io.EOF {
		response.Failed(w, err)
		return
	}

	response.Success(w, "check ok")
}

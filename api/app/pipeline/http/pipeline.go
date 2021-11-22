package http

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/infraboard/keyauth/app/token"
	"github.com/infraboard/mcube/grpc/gcontext"
	"github.com/infraboard/mcube/http/context"
	"github.com/infraboard/mcube/http/request"
	"github.com/infraboard/mcube/http/response"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/infraboard/workflow/api/app/pipeline"
)

func (h *handler) CreatePipeline(w http.ResponseWriter, r *http.Request) {
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

	req := pipeline.NewCreatePipelineRequest()
	if err := request.GetDataFromRequest(r, req); err != nil {
		response.Failed(w, err)
		return
	}
	req.UpdateOwner(tk)

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
	tk, ok := hc.AuthInfo.(*token.Token)
	if !ok {
		response.Failed(w, fmt.Errorf("auth info is not an *token.Token"))
		return
	}

	req := pipeline.NewDescribePipelineRequestWithID(hc.PS.ByName("id"))
	req.Namespace = tk.Namespace

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

var (
	// 升级为ws协议
	upgrader = websocket.Upgrader{
		HandshakeTimeout: 60 * time.Second,
		ReadBufferSize:   8192,
		WriteBufferSize:  8192,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func (h *handler) WatchPipeline(w http.ResponseWriter, r *http.Request) {
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

	var header, trailer metadata.MD
	rpcCtx := ctx.Context()
	stream, err := h.service.WatchPipeline(
		rpcCtx,
		grpc.Header(&header),
		grpc.Trailer(&trailer),
	)
	defer func() {
		rpcCtx.Done()
	}()

	if err != nil {
		response.Failed(w, gcontext.NewExceptionFromTrailer(trailer, err))
		return
	}

	err = stream.Send(&pipeline.WatchPipelineRequest{
		RequestUnion: &pipeline.WatchPipelineRequest_CreateRequest{CreateRequest: req},
	})

	if err != nil {
		h.log.Errorf("stream send watch req error, %s", err)
		return
	}

	var responseHeader http.Header
	// If Sec-WebSocket-Protocol starts with "Bearer", respond in kind.
	// TODO(tmc): consider customizability/extension point here.
	if strings.HasPrefix(r.Header.Get("Sec-WebSocket-Protocol"), "Bearer") {
		responseHeader = http.Header{
			"Sec-WebSocket-Protocol": []string{"Bearer"},
		}
	}

	// https --> websocket
	conn, err := upgrader.Upgrade(w, r, responseHeader)
	if err != nil {
		h.log.Errorf("error upgrading websocket:", err)
		return
	}
	defer conn.Close()

	dumpper := NewPipelineStreamDumpper(stream)
	defer dumpper.Close()
	h.proxy.Proxy(r.Context(), conn, dumpper)
}

func NewPipelineStreamDumpper(stream pipeline.Service_WatchPipelineClient) *PipelineStreamDumpper {
	return &PipelineStreamDumpper{
		stream: stream,
	}
}

type PipelineStreamDumpper struct {
	stream  pipeline.Service_WatchPipelineClient
	watchId int64
}

func (d *PipelineStreamDumpper) Read(buf []byte) (n int, err error) {
	pp, err := d.stream.Recv()
	if err != nil {
		return 0, err
	}
	d.watchId = pp.WatchId

	if pp.Pipeline == nil {
		return 0, nil
	}

	data, err := json.Marshal(pp.Pipeline)
	if err != nil {
		return 0, err
	}

	buf = append(buf, data...)
	fmt.Println(string(buf), len(buf))
	return len(data), nil
}

func (d *PipelineStreamDumpper) Write(buf []byte) (n int, err error) {
	return 0, nil
}

func (d *PipelineStreamDumpper) Close() error {
	cancelReq := &pipeline.CancelWatchPipelineRequest{WatchId: d.watchId}

	req := &pipeline.WatchPipelineRequest{
		RequestUnion: &pipeline.WatchPipelineRequest_CancelRequest{CancelRequest: cancelReq},
	}

	return d.stream.Send(req)
}

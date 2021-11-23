package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/infraboard/keyauth/app/token"
	"github.com/infraboard/mcube/http/context"
	"github.com/infraboard/mcube/http/request"
	"github.com/infraboard/mcube/http/response"

	"github.com/infraboard/workflow/api/app/pipeline"
)

func (h *handler) CreatePipeline(w http.ResponseWriter, r *http.Request) {
	ctx := context.GetContext(r)
	tk := ctx.AuthInfo.(*token.Token)

	req := pipeline.NewCreatePipelineRequest()
	if err := request.GetDataFromRequest(r, req); err != nil {
		response.Failed(w, err)
		return
	}
	req.UpdateOwner(tk)

	ins, err := h.service.CreatePipeline(
		r.Context(),
		req,
	)
	if err != nil {
		response.Failed(w, err)
		return
	}
	response.Success(w, ins)
}

func (h *handler) QueryPipeline(w http.ResponseWriter, r *http.Request) {
	page := request.NewPageRequestFromHTTP(r)
	req := pipeline.NewQueryPipelineRequest()
	req.Page = &page.PageRequest

	dommains, err := h.service.QueryPipeline(
		r.Context(),
		req,
	)
	if err != nil {
		response.Failed(w, err)
		return
	}
	response.Success(w, dommains)
}

func (h *handler) DescribePipeline(w http.ResponseWriter, r *http.Request) {
	ctx := context.GetContext(r)
	tk := ctx.AuthInfo.(*token.Token)

	req := pipeline.NewDescribePipelineRequestWithID(ctx.PS.ByName("id"))
	req.Namespace = tk.Namespace

	dommains, err := h.service.DescribePipeline(
		r.Context(),
		req,
	)
	if err != nil {
		response.Failed(w, err)
		return
	}
	response.Success(w, dommains)
}

// pipeline删除时,除了删除pipeline对象本身而外，还需要删除该pipeline下的所有step
func (h *handler) DeletePipeline(w http.ResponseWriter, r *http.Request) {
	ctx := context.GetContext(r)
	req := pipeline.NewDeletePipelineRequestWithID(ctx.PS.ByName("id"))

	dommains, err := h.service.DeletePipeline(
		r.Context(),
		req,
	)
	if err != nil {
		response.Failed(w, err)
		return
	}
	response.Success(w, dommains)
}

func (h *handler) WatchPipelineCheck(w http.ResponseWriter, r *http.Request) {
	ctx := context.GetContext(r)
	tk := ctx.AuthInfo.(*token.Token)

	req := pipeline.NewWatchPipelineRequestByID("", ctx.PS.ByName("id"))
	req.Namespace = tk.Namespace
	req.DryRun = true

	err := h.service.WatchPipeline(nil)

	if err != nil {
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
	ctx := context.GetContext(r)
	tk := ctx.AuthInfo.(*token.Token)

	req := pipeline.NewWatchPipelineRequestByID("", ctx.PS.ByName("id"))
	req.Namespace = tk.Namespace

	err := h.service.WatchPipeline(nil)
	if err != nil {
		response.Failed(w, err)
		return
	}

	// err = stream.Send(&pipeline.WatchPipelineRequest{
	// 	RequestUnion: &pipeline.WatchPipelineRequest_CreateRequest{CreateRequest: req},
	// })

	// if err != nil {
	// 	h.log.Errorf("stream send watch req error, %s", err)
	// 	return
	// }

	// var responseHeader http.Header
	// // If Sec-WebSocket-Protocol starts with "Bearer", respond in kind.
	// // TODO(tmc): consider customizability/extension point here.
	// if strings.HasPrefix(r.Header.Get("Sec-WebSocket-Protocol"), "Bearer") {
	// 	responseHeader = http.Header{
	// 		"Sec-WebSocket-Protocol": []string{"Bearer"},
	// 	}
	// }

	// // https --> websocket
	// conn, err := upgrader.Upgrade(w, r, responseHeader)
	// if err != nil {
	// 	h.log.Errorf("error upgrading websocket:", err)
	// 	return
	// }
	// defer conn.Close()

	// dumpper := NewPipelineStreamDumpper(stream)
	// defer dumpper.Close()
	// h.proxy.Proxy(r.Context(), conn, dumpper)
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

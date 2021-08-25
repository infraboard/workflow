package http

import (
	"net/http"

	"github.com/infraboard/mcube/grpc/gcontext"
	"github.com/infraboard/mcube/http/context"
	"github.com/infraboard/mcube/http/request"
	"github.com/infraboard/mcube/http/response"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/infraboard/workflow/api/pkg/pipeline"
)

// Action
func (h *handler) CreateStep(w http.ResponseWriter, r *http.Request) {
	ctx, err := gcontext.NewGrpcOutCtxFromHTTPRequest(r)
	if err != nil {
		response.Failed(w, err)
		return
	}

	req := pipeline.NewCreateStepRequest()
	if err := request.GetDataFromRequest(r, req); err != nil {
		response.Failed(w, err)
		return
	}

	var header, trailer metadata.MD
	ins, err := h.service.CreateStep(
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

func (h *handler) QueryStep(w http.ResponseWriter, r *http.Request) {
	ctx, err := gcontext.NewGrpcOutCtxFromHTTPRequest(r)
	if err != nil {
		response.Failed(w, err)
		return
	}

	page := request.NewPageRequestFromHTTP(r)
	req := pipeline.NewQueryStepRequest()
	req.Page = &page.PageRequest

	var header, trailer metadata.MD
	dommains, err := h.service.QueryStep(
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

func (h *handler) DescribeStep(w http.ResponseWriter, r *http.Request) {
	ctx, err := gcontext.NewGrpcOutCtxFromHTTPRequest(r)
	if err != nil {
		response.Failed(w, err)
		return
	}

	hc := context.GetContext(r)
	req := pipeline.NewDescribeStepRequestWithKey(hc.PS.ByName("id"))

	var header, trailer metadata.MD
	dommains, err := h.service.DescribeStep(
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

func (h *handler) DeleteStep(w http.ResponseWriter, r *http.Request) {
	ctx, err := gcontext.NewGrpcOutCtxFromHTTPRequest(r)
	if err != nil {
		response.Failed(w, err)
		return
	}

	hc := context.GetContext(r)
	req := pipeline.NewDeleteStepRequestWithKey(hc.PS.ByName("id"))

	var header, trailer metadata.MD
	dommains, err := h.service.DeleteStep(
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

	var header, trailer metadata.MD
	dommains, err := h.service.AuditStep(
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

func (h *handler) CancelStep(w http.ResponseWriter, r *http.Request) {
	ctx, err := gcontext.NewGrpcOutCtxFromHTTPRequest(r)
	if err != nil {
		response.Failed(w, err)
		return
	}

	hc := context.GetContext(r)
	req := pipeline.NewCancelStepRequestWithKey(hc.PS.ByName("id"))

	var header, trailer metadata.MD
	dommains, err := h.service.CancelStep(
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

func (h *handler) QueryVariableTemplate(w http.ResponseWriter, r *http.Request) {
	if !tempateIsInit {
		for k, v := range pipeline.VALUE_TYPE_ID_MAP {
			for i := range valueTempate {
				if valueTempate[i].Type == v {
					valueTempate[i].Prefix = k
				}
			}
		}
		tempateIsInit = true
	}

	response.Success(w, valueTempate)
}

var (
	tempateIsInit = false
	valueTempate  = []*ValueTypeDesc{
		{
			Type:   pipeline.PARAM_VALUE_TYPE_PLAIN,
			Prefix: "",
			Name:   "明文",
			Desc:   "明文文本,敏感信息请不要使用这个类型",
			IsEdit: true,
		},
		{
			Type:   pipeline.PARAM_VALUE_TYPE_PASSWORD,
			Prefix: "",
			Name:   "秘文",
			Desc:   "敏感信息,由系统加密存储,运行时解密注入",
			IsEdit: true,
		},
		{
			Type:   pipeline.PARAM_VALUE_TYPE_CRYPTO,
			Prefix: "",
			Name:   "解密",
			Desc:   "敏感信息加密后的密文,无法修改",
		},
		{
			Type:   pipeline.PARAM_VALUE_TYPE_APP_VAR,
			Prefix: "",
			Name:   "应用变量",
			Desc:   "应用属性,也包含自定义变量,运行时由系统动态注入",
			IsEdit: true,
		},
		{
			Type:   pipeline.PARAM_VALUE_TYPE_SECRET_REF,
			Prefix: "",
			Name:   "Secret引用",
			Desc:   "运行时由系统查询Secret系统后动态注入",
			IsEdit: true,
		},
	}
)

type ValueTypeDesc struct {
	Type   pipeline.PARAM_VALUE_TYPE `json:"type"`
	Prefix string                    `json:"prefix"`
	Name   string                    `json:"name"`
	Desc   string                    `json:"desc"`
	Value  string                    `json:"value"`
	IsEdit bool                      `json:"is_edit"`
}

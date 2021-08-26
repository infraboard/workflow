package http

import (
	"net/http"

	"github.com/infraboard/mcube/http/response"

	"github.com/infraboard/workflow/api/pkg/pipeline"
	"github.com/infraboard/workflow/common/enum"
)

func (h *handler) QueryStepStatusEnum(w http.ResponseWriter, r *http.Request) {
	response.Success(w, stepStatusEnums)
}

var (
	stepStatusEnums = []*enum.EnumDesc{
		{Value: pipeline.STEP_STATUS_PENDDING.String(), Name: "调度中", Desc: "当任务被创建后触发"},
		{Value: pipeline.STEP_STATUS_RUNNING.String(), Name: "运行中", Desc: "当任务开始运行时触发"},
		{Value: pipeline.STEP_STATUS_SUCCEEDED.String(), Name: "运行成功", Desc: "当任务运行成功后触发"},
		{Value: pipeline.STEP_STATUS_FAILED.String(), Name: "运行失败", Desc: "当任务运行失败时触发"},
		{Value: pipeline.STEP_STATUS_CANCELED.String(), Name: "任务取消", Desc: "当任务被成功取消时触发"},
		{Value: pipeline.STEP_STATUS_SKIP.String(), Name: "任务跳过", Desc: "当前任务忽略执行时触发"},
		{Value: pipeline.STEP_STATUS_AUDITING.String(), Name: "审批中", Desc: "当前任务需要审批时触发"},
		{Value: pipeline.STEP_STATUS_REFUSE.String(), Name: "审批拒绝", Desc: "当审批的任务被拒绝时触发"},
		{Value: pipeline.STEP_STATUS_SCHEDULE_FAILED.String(), Name: "调度失败", Desc: "当任务由于调度失败无法执行时触发"},
	}
)

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

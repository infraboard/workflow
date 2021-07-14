package pipeline

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/infraboard/keyauth/pkg/token"
	"github.com/infraboard/mcube/http/request"
	"github.com/infraboard/mcube/pb/resource"
	"github.com/infraboard/mcube/types/ftime"
)

func NewCreateActionRequest() *CreateActionRequest {
	return &CreateActionRequest{}
}

// NewQueryActionRequest 查询book列表
func NewQueryActionRequest(page *request.PageRequest) *QueryActionRequest {
	return &QueryActionRequest{
		Page: &page.PageRequest,
	}
}

func (req *CreateActionRequest) Validate() error {
	return validate.Struct(req)
}

func LoadActionFromBytes(payload []byte) (*Action, error) {
	ins := NewDefaultAction()

	// 解析Value
	if err := json.Unmarshal(payload, ins); err != nil {
		return nil, fmt.Errorf("unmarshal step error, vaule(%s) %s", string(payload), err)
	}

	// 校验合法性
	if err := ins.Validate(); err != nil {
		return nil, err
	}

	return ins, nil
}

func NewDefaultAction() *Action {
	return &Action{}
}

func NewAction(req *CreateActionRequest) (*Action, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	p := &Action{
		CreateAt:     ftime.Now().Timestamp(),
		UpdateAt:     ftime.Now().Timestamp(),
		Name:         req.Name,
		VisiableMode: req.VisiableMode,
		RunnerType:   req.RunnerType,
		RunnerParams: req.RunnerParams,
		RunParams:    req.RunParams,
		Tags:         req.Tags,
		Description:  req.Description,
	}

	return p, nil
}

func (a *Action) UpdateOwner(tk *token.Token) {
	a.CreateBy = tk.Account
	a.Domain = tk.Domain
	a.Namespace = tk.Namespace
}

func (a *Action) DefaultRunParam() map[string]string {
	param := map[string]string{}
	for k, v := range a.RunParams {
		if v != nil && v.Default != "" {
			param[k] = v.Default
		}
	}
	return param
}

func (a *Action) DefaultRunnerParam() map[string]string {
	param := map[string]string{}
	for k, v := range a.RunnerParams {
		if v != nil && v.Default != "" {
			param[k] = v.Default
		}
	}
	return param
}

// ValidateParam 按照action的定义, 检查必传参数是否传人
func (a *Action) ValidateRunParam(params map[string]string) error {
	msg := []string{}
	for k, v := range a.RunParams {
		if v != nil && v.Required {
			if pv, ok := params[k]; !ok || pv == "" {
				msg = append(msg, "required param %s", k)
			}
		}
	}

	if len(msg) > 0 {
		return fmt.Errorf("validate run params error, %s", strings.Join(msg, ","))
	}

	return nil
}

// ValidateParam 按照action的定义, 检查必传参数是否传人
func (a *Action) ValidateRunnerParam(params map[string]string) error {
	msg := []string{}
	for k, v := range a.RunnerParams {
		if v != nil && v.Required {
			if pv, ok := params[k]; !ok || pv == "" {
				msg = append(msg, "required param %s", k)
			}
		}
	}

	if len(msg) > 0 {
		return fmt.Errorf("validate runner params error, %s", strings.Join(msg, ","))
	}

	return nil
}

func (a *Action) Validate() error {
	return validate.Struct(a)
}

func (a *Action) EtcdObjectKey() string {
	ns := a.VisiableMode.String()
	if a.VisiableMode.Equal(resource.VisiableMode_NAMESPACE) {
		ns = a.Namespace
	}
	return fmt.Sprintf("%s/%s/%s", EtcdActionPrefix(), ns, a.Name)
}

// NewActionSet todo
func NewActionSet() *ActionSet {
	return &ActionSet{
		Items: []*Action{},
	}
}

func (s *ActionSet) Add(item *Action) {
	s.Items = append(s.Items, item)
}

// NewQueryActionRequest 查询book列表
func NewDescribeActionRequestWithName(name string) *DescribeActionRequest {
	return &DescribeActionRequest{
		Name: name,
	}
}

// NewDeleteActionRequestWithID 查询book列表
func NewDeleteActionRequestWithName(name string) *DeleteActionRequest {
	return &DeleteActionRequest{
		Name: name,
	}
}

func (req *DeleteActionRequest) Namespace(tk *token.Token) string {
	ns := req.VisiableMode.String()
	if req.VisiableMode.Equal(resource.VisiableMode_NAMESPACE) {
		ns = tk.Namespace
	}

	return ns
}

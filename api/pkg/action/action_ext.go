package action

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/infraboard/keyauth/pkg/token"
	"github.com/infraboard/mcube/http/request"
	"github.com/infraboard/mcube/pb/resource"
	"github.com/infraboard/mcube/types/ftime"
	"github.com/rs/xid"
)

// use a single instance of Validate, it caches struct info
var (
	validate = validator.New()
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
		Id:           xid.New().String(),
		CreateAt:     ftime.Now().Timestamp(),
		UpdateAt:     ftime.Now().Timestamp(),
		Name:         req.Name,
		Version:      req.Version,
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
	ret := map[string]string{}
	for i := range a.RunParams {
		param := a.RunParams[i]
		if param.DefaultValue != "" {
			ret[param.KeyName] = param.DefaultValue
		}
	}
	return ret
}

func (a *Action) RunnerParam() map[string]string {
	param := map[string]string{}
	for k, v := range a.RunnerParams {
		if v != nil && v.Value != "" {
			param[k] = v.Value
		}
	}
	return param
}

// ValidateParam 按照action的定义, 检查必传参数是否传人
func (a *Action) ValidateRunParam(params map[string]string) error {
	msg := []string{}
	for i := range a.RunParams {
		param := a.RunParams[i]
		if param.Required {
			if pv, ok := params[param.KeyName]; !ok || pv == "" {
				msg = append(msg, "required param %s", pv)
			}
		}
	}

	if len(msg) > 0 {
		return fmt.Errorf("validate run params error, %s", strings.Join(msg, ","))
	}

	return nil
}

func (a *Action) Validate() error {
	return validate.Struct(a)
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
func NewDescribeActionRequest(namespace, name, version string) *DescribeActionRequest {
	return &DescribeActionRequest{
		Namespace: namespace,
		Name:      name,
		Version:   version,
	}
}

// NewDeleteActionRequest 查询book列表
func NewDeleteActionRequest(version, name string) *DeleteActionRequest {
	return &DeleteActionRequest{
		Name:    name,
		Version: version,
	}
}

func (req *DeleteActionRequest) Namespace(tk *token.Token) string {
	ns := req.VisiableMode.String()
	if req.VisiableMode.Equal(resource.VisiableMode_NAMESPACE) {
		ns = tk.Namespace
	}

	return ns
}

func (req *DescribeActionRequest) Validate() error {
	return validate.Struct(req)
}

func ParseActionKey(key string) (name, version string) {
	parseKey := strings.Split(key, "@")
	name = parseKey[0]
	if len(parseKey) > 1 {
		version = parseKey[1]
	}

	return
}

package action

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/infraboard/keyauth/app/token"
	"github.com/infraboard/mcube/http/request"
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
		Page: page,
	}
}

func (req *CreateActionRequest) Validate() error {
	return validate.Struct(req)
}

func (req *CreateActionRequest) UpdateOwner(tk *token.Token) {
	req.CreateBy = tk.Account
	req.Domain = tk.Domain
	req.Namespace = tk.Namespace
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
		CreateAt:     time.Now().UnixMilli(),
		UpdateAt:     time.Now().UnixMilli(),
		Domain:       req.Domain,
		Namespace:    req.Namespace,
		CreateBy:     req.CreateBy,
		Logo:         req.Logo,
		DisplayName:  req.DisplayName,
		IsLatest:     true,
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

func (a *Action) InitNil() {
	if a.RunnerParams == nil {
		a.RunnerParams = map[string]string{}
	}
	if a.RunParams == nil {
		a.RunParams = []*RunParamDesc{}
	}
	if a.Tags == nil {
		a.Tags = map[string]string{}
	}
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
		if v != "" {
			param[k] = v
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
				msg = append(msg, "required param "+param.KeyName)
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

func (a *Action) Update(req *UpdateActionRequest) {
	a.VisiableMode = req.VisiableMode
	a.RunParams = req.RunParams
	a.Tags = req.Tags
	a.Description = req.Description
}

func (a *Action) Key() string {
	return fmt.Sprintf("%s@%s", a.Name, a.Version)
}

// NewActionSet todo
func NewActionSet() *ActionSet {
	return &ActionSet{
		Items: []*Action{},
	}
}

func (s *ActionSet) InitNil() {
	for i := range s.Items {
		s.Items[i].InitNil()
	}
}

func (s *ActionSet) Add(item *Action) {
	s.Items = append(s.Items, item)
}

// NewQueryActionRequest 查询book列表
func NewDescribeActionRequest(name, version string) *DescribeActionRequest {
	return &DescribeActionRequest{
		Name:    name,
		Version: version,
	}
}

// NewDeleteActionRequest 查询book列表
func NewDeleteActionRequest(name, version string) *DeleteActionRequest {
	return &DeleteActionRequest{
		Name:    name,
		Version: version,
	}
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

func (req *DeleteActionRequest) Validate() error {
	return validate.Struct(req)
}

func NewUpdateActionRequest() *UpdateActionRequest {
	return &UpdateActionRequest{}
}

func (req *UpdateActionRequest) Validate() error {
	return validate.Struct(req)
}

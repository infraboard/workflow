package pipeline

import (
	"encoding/json"
	"fmt"

	"github.com/infraboard/mcube/http/request"
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
		RunParams:    req.RunParams,
		Tags:         req.Tags,
		Description:  req.Description,
	}

	return p, nil
}

func (p *Action) Validate() error {
	return validate.Struct(p)
}

func (p *Action) EtcdObjectKey() string {
	return fmt.Sprintf("%s/%s", EtcdActionPrefix(), p.Name)

	// return fmt.Sprintf("%s/%s/%s", EtcdActionPrefix(), p.Namespace, p.Name)
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

package template

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/infraboard/keyauth/app/token"
	"github.com/infraboard/mcube/http/request"
	"github.com/infraboard/mcube/types/ftime"
	"github.com/rs/xid"
)

// use a single instance of Validate, it caches struct info
var (
	validate = validator.New()
)

func NewTemplate(req *CreateTemplateRequest) (*Template, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	p := &Template{
		Domain:       req.Domain,
		Namespace:    req.Namespace,
		CreateBy:     req.CreateBy,
		Id:           xid.New().String(),
		CreateAt:     ftime.Now().Timestamp(),
		UpdateAt:     ftime.Now().Timestamp(),
		Name:         req.Name,
		Tags:         req.Tags,
		VisiableMode: req.VisiableMode,
		Pipelines:    req.Pipelines,
		Description:  req.Description,
	}

	return p, nil
}

func (req *CreateTemplateRequest) Validate() error {
	// 判断同一个模版里面,Pipeline名字是否相同
	nameMap := map[string]struct{}{}
	for i := range req.Pipelines {
		_, ok := nameMap[req.Pipelines[i].Name]
		if ok {
			return fmt.Errorf("in this template name is ready exist")
		}
		nameMap[req.Pipelines[i].Name] = struct{}{}
	}

	return validate.Struct(req)
}

func NewCreateTemplateRequest() *CreateTemplateRequest {
	return &CreateTemplateRequest{}
}

func (req *CreateTemplateRequest) UpdateOwner(tk *token.Token) {
	req.Domain = tk.Domain
	req.Namespace = tk.Namespace
	req.CreateBy = tk.Account
}

// NewTemplateSet todo
func NewTemplateSet() *TemplateSet {
	return &TemplateSet{
		Items: []*Template{},
	}
}

func (s *TemplateSet) Add(item *Template) {
	s.Items = append(s.Items, item)
}

func NewDefaultTemplate() *Template {
	return &Template{}
}

func (t *Template) Update(updater string, req *UpdateTemplateData) {
	t.UpdateAt = ftime.Now().Timestamp()
	t.UpdateBy = updater
	t.Name = req.Name
	t.Tags = req.Tags
	t.Description = req.Description
	t.Pipelines = req.Pipelines
}

func (t *Template) Patch(updater string, req *UpdateTemplateData) {
	t.UpdateAt = ftime.Now().Timestamp()
	t.UpdateBy = updater

	if req.Name != "" {
		t.Name = req.Name
	}
	if req.Description != "" {
		t.Description = req.Description
	}
	if len(req.Tags) > 0 {
		t.Tags = req.Tags
	}
	if len(req.Pipelines) > 0 {
		t.Pipelines = req.Pipelines
	}
}

func (req *DescribeTemplateRequest) Validate() error {
	return validate.Struct(req)
}

// NewQueryTemplateRequest 查询book列表
func NewQueryTemplateRequest(page *request.PageRequest) *QueryTemplateRequest {
	return &QueryTemplateRequest{
		Page: &page.PageRequest,
	}
}

// NewDescribeTemplateRequestWithID 查询book列表
func NewDescribeTemplateRequestWithID(id string) *DescribeTemplateRequest {
	return &DescribeTemplateRequest{
		Id: id,
	}
}

// NewDeleteTemplateRequestWithID 查询book列表
func NewDeleteTemplateRequestWithID(id string) *DeleteTemplateRequest {
	return &DeleteTemplateRequest{
		Id: id,
	}
}

func (req *UpdateTemplateRequest) Validate() error {
	return validate.Struct(req)
}

func NewUpdateTemplateRequest(id string) *UpdateTemplateRequest {
	return &UpdateTemplateRequest{
		Id:   id,
		Data: &UpdateTemplateData{},
	}
}

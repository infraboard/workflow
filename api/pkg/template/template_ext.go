package template

import (
	"github.com/go-playground/validator/v10"
	"github.com/infraboard/keyauth/pkg/token"
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
		Pipeline:     req.Pipeline,
		Description:  req.Description,
	}

	return p, nil
}

func (req *CreateTemplateRequest) Validate() error {
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

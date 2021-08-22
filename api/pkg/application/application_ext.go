package application

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/infraboard/keyauth/pkg/token"
	"github.com/infraboard/mcube/http/request"
	"github.com/infraboard/mcube/types/ftime"
	"github.com/infraboard/workflow/common/repo/gitlab"
	"github.com/rs/xid"
)

// use a single instance of Validate, it caches struct info
var (
	validate = validator.New()
)

func NewDefaultApplication() *Application {
	return &Application{}
}

func NewCreateApplicationRequest() *CreateApplicationRequest {
	return &CreateApplicationRequest{}
}

func (req *CreateApplicationRequest) UpdateOwner(tk *token.Token) {
	req.CreateBy = tk.Account
	req.Domain = tk.Domain
	req.Namespace = tk.Namespace
}

func (req *CreateApplicationRequest) NeedSetHook() bool {
	return req.ScmProjectId != "" && req.ScmPrivateToken != ""
}

func (req *CreateApplicationRequest) GetScmAddr() (string, error) {
	if req.RepoHttpUrl == "" {
		return "", nil
	}

	url, err := url.Parse(req.RepoHttpUrl)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s://%s", url.Scheme, url.Host), nil
}

func (req *CreateApplicationRequest) Int64ScmProjectID() int64 {
	id, _ := strconv.ParseInt(req.ScmProjectId, 10, 64)
	return id
}

func (req *CreateApplicationRequest) Validate() error {
	return validate.Struct(req)
}

// NewApplication todo
func NewApplication(req *CreateApplicationRequest) (*Application, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	ins := &Application{
		Id:           xid.New().String(),
		CreateAt:     ftime.Now().Timestamp(),
		UpdateAt:     ftime.Now().Timestamp(),
		Domain:       req.Domain,
		Namespace:    req.Namespace,
		CreateBy:     req.CreateBy,
		Name:         req.Name,
		Tags:         req.Tags,
		Description:  req.Description,
		Pipeline:     req.Pipeline,
		RepoSshUrl:   req.RepoSshUrl,
		RepoHttpUrl:  req.RepoHttpUrl,
		ScmType:      req.ScmType,
		ScmProjectId: req.ScmProjectId,
	}

	return ins, nil
}

func (a *Application) AddError(err error) {
	a.Errors = append(a.Errors, err.Error())
}

func (a *Application) GenWebHook(callbackURL string) *gitlab.WebHook {
	return &gitlab.WebHook{
		PushEvents:          true,
		TagPushEvents:       true,
		MergeRequestsEvents: true,
		Token:               a.Id,
		Url:                 callbackURL,
	}
}

// NewApplicationSet 实例
func NewApplicationSet() *ApplicationSet {
	return &ApplicationSet{
		Items: []*Application{},
	}
}

func (s *ApplicationSet) Add(item *Application) {
	s.Items = append(s.Items, item)
}

// NewQueryBookRequest 查询book列表
func NewQueryBookRequest(page *request.PageRequest) *QueryApplicationRequest {
	return &QueryApplicationRequest{
		Page: &page.PageRequest,
	}
}

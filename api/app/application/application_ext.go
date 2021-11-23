package application

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/rs/xid"
	"google.golang.org/protobuf/proto"

	"github.com/infraboard/keyauth/app/token"
	"github.com/infraboard/mcube/http/request"
	pb "github.com/infraboard/mcube/pb/request"
	"github.com/infraboard/mcube/types/ftime"

	"github.com/infraboard/workflow/api/app/pipeline"
	"github.com/infraboard/workflow/api/app/scm"
	"github.com/infraboard/workflow/api/app/scm/gitlab"
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
	names := map[string]struct{}{}
	for i := range req.Pipelines {
		p := req.Pipelines[i]
		_, ok := names[p.Name]
		if ok {
			return fmt.Errorf("pipeline name %s ready exist", p.Name)
		}
		names[p.Name] = struct{}{}
	}

	return validate.Struct(req)
}

// NewApplication todo
func NewApplication(req *CreateApplicationRequest) (*Application, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	ins := &Application{
		Id:          xid.New().String(),
		CreateAt:    ftime.Now().Timestamp(),
		UpdateAt:    ftime.Now().Timestamp(),
		Domain:      req.Domain,
		Namespace:   req.Namespace,
		CreateBy:    req.CreateBy,
		Name:        req.Name,
		Tags:        req.Tags,
		Description: req.Description,
		// Pipelines:       req.Pipelines,
		RepoSshUrl:      req.RepoSshUrl,
		RepoHttpUrl:     req.RepoHttpUrl,
		ScmType:         req.ScmType,
		ScmProjectId:    req.ScmProjectId,
		ScmPrivateToken: req.ScmPrivateToken,
	}

	return ins, nil
}

func (a *Application) Update(updater string, req *UpdateApplicationData) {
	a.UpdateAt = ftime.Now().Timestamp()
	a.UpdateBy = updater
	a.Name = req.Name
	a.Tags = req.Tags
	a.Description = req.Description
	// a.Pipelines = req.Pipelines
}

func (a *Application) Patch(updater string, req *UpdateApplicationData) {
	a.UpdateAt = ftime.Now().Timestamp()
	a.UpdateBy = updater

	if req.Name != "" {
		a.Name = req.Name
	}
	if req.Description != "" {
		a.Description = req.Description
	}
	if len(req.Tags) > 0 {
		a.Tags = req.Tags
	}
	if len(req.Pipelines) > 0 {
		// a.Pipelines = req.Pipelines
	}
}

func (a *Application) Desensitize() {
	a.ScmPrivateToken = "****"
}

func (a *Application) GenWebHook(callbackURL string) *gitlab.WebHook {
	cb := fmt.Sprintf("%s/workflow/api/v1/triggers/scm/%s",
		callbackURL, strings.ToLower(a.ScmType.String()))

	return &gitlab.WebHook{
		PushEvents:          true,
		TagPushEvents:       true,
		MergeRequestsEvents: true,
		Token:               a.Id,
		Url:                 cb,
	}
}

func (a *Application) GetScmAddr() (string, error) {
	if a.RepoHttpUrl == "" {
		return "", nil
	}

	url, err := url.Parse(a.RepoHttpUrl)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s://%s", url.Scheme, url.Host), nil
}

func (a *Application) Int64ScmProjectID() int64 {
	id, _ := strconv.ParseInt(a.ScmProjectId, 10, 64)
	return id
}

func (a *Application) Int64ScmHookID() int64 {
	id, _ := strconv.ParseInt(a.ScmHookId, 10, 64)
	return id
}

func (a *Application) MatchPipeline(e *scm.WebHookEvent) (mached []*pipeline.CreatePipelineRequest) {
	for i := range a.Pipelines {
		if a.Pipelines[i].Define.On.IsMatch(e.GetBranche(), e.GetEventName()) {
			mached = append(mached, a.Pipelines[i].Define)
		}
	}
	return
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

// NewQueryApplicationRequest 查询book列表
func NewQueryApplicationRequest(page *request.PageRequest) *QueryApplicationRequest {
	return &QueryApplicationRequest{
		Page: &page.PageRequest,
	}
}

func NewDescribeApplicationRequestWithID(id string) *DescribeApplicationRequest {
	return &DescribeApplicationRequest{
		Id: id,
	}
}

func NewDescribeApplicationRequestWithName(namespace, name string) *DescribeApplicationRequest {
	return &DescribeApplicationRequest{
		Namespace: namespace,
		Name:      name,
	}
}

func (r *DescribeApplicationRequest) Validate() error {
	return validate.Struct(r)
}

func NewDeleteApplicationRequest(namespace, name string) *DeleteApplicationRequest {
	return &DeleteApplicationRequest{
		Namespace: namespace,
		Name:      name,
	}
}

func (req *ApplicationEvent) Validate() error {
	if req.WebhookEvent == nil {
		return fmt.Errorf("web hook event is nil")
	}

	return validate.Struct(req)
}

func NewPipelineCreateStatus() *PipeLineStatus {
	return &PipeLineStatus{
		Timestamp: ftime.Now().Timestamp(),
	}
}

func NewApplicationEvent(appid string, event *scm.WebHookEvent) *ApplicationEvent {
	return &ApplicationEvent{
		AppId:        appid,
		WebhookEvent: event,
	}
}

func (s *PipeLineStatus) Clone() *PipeLineStatus {
	return proto.Clone(s).(*PipeLineStatus)
}

func NewUpdateApplicationRequest(appId string) *UpdateApplicationRequest {
	return &UpdateApplicationRequest{
		Id:         appId,
		UpdateMode: pb.UpdateMode_PUT,
		Data:       &UpdateApplicationData{},
	}
}

func (req *UpdateApplicationRequest) Validate() error {
	return validate.Struct(req)
}

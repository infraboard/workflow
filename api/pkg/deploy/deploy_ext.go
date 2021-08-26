package deploy

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/infraboard/mcube/http/request"
	"github.com/infraboard/mcube/types/ftime"
	"github.com/rs/xid"

	"github.com/infraboard/keyauth/pkg/token"
)

// use a single instance of Validate, it caches struct info
var (
	validate = validator.New()
)

func NewDefaultApplicationDeploy() *ApplicationDeploy {
	return &ApplicationDeploy{}
}

// NewApplication todo
func NewApplicationDeploy(req *CreateApplicationDeployRequest) (*ApplicationDeploy, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	ins := &ApplicationDeploy{
		Id:               xid.New().String(),
		CreateAt:         ftime.Now().Timestamp(),
		UpdateAt:         ftime.Now().Timestamp(),
		Domain:           req.Domain,
		Namespace:        req.Namespace,
		CreateBy:         req.CreateBy,
		AppId:            req.AppId,
		Environment:      req.Environment,
		Name:             req.Name,
		ApplicationConf:  req.ApplicationConf,
		DeployMode:       req.DeployMode,
		HostDeployConfig: req.HostDeployConfig,
		K8SDeployConfig:  req.K8SDeployConfig,
		Tags:             req.Tags,
		Description:      req.Description,
	}

	return ins, nil
}

func (a *ApplicationDeploy) Desensitize() {
}

func (req *CreateApplicationDeployRequest) Validate() error {
	switch req.DeployMode {
	case Mode_Host:
		if req.HostDeployConfig == nil {
			return fmt.Errorf("host mode but HostDeployConfig is nil")
		}
	case Mode_K8s:
		if req.K8SDeployConfig == nil {
			return fmt.Errorf("k8s mode but K8SDeployConfig is nil")
		}
	default:
		return fmt.Errorf("unknown deploy type %s", req.DeployMode)
	}

	return validate.Struct(req)
}

// NewApplicationSet 实例
func NewApplicationDeploySet() *ApplicationDeploySet {
	return &ApplicationDeploySet{
		Items: []*ApplicationDeploy{},
	}
}

func (s *ApplicationDeploySet) Add(item *ApplicationDeploy) {
	s.Items = append(s.Items, item)
}

func NewCreateApplicationDeployRequest() *CreateApplicationDeployRequest {
	return &CreateApplicationDeployRequest{}
}

func (req *CreateApplicationDeployRequest) UpdateOwner(tk *token.Token) {
	req.CreateBy = tk.Account
	req.Domain = tk.Domain
	req.Namespace = tk.Namespace
}

// NewQueryApplicationDeployRequest 查询book列表
func NewQueryApplicationDeployRequest(page *request.PageRequest) *QueryApplicationDeployRequest {
	return &QueryApplicationDeployRequest{
		Page: &page.PageRequest,
	}
}

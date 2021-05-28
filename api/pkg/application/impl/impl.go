package impl

import (
	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/logger/zap"
	"github.com/infraboard/mcube/pb/http"

	"github.com/infraboard/workflow/api/pkg"
	"github.com/infraboard/workflow/api/pkg/application"
)

var (
	// Service 服务实例
	Service = &service{}
)

type service struct {
	application.UnimplementedServiceServer

	log logger.Logger
}

func (s *service) Config() error {
	// get global config with here
	s.log = zap.L().Named("Application")
	return nil
}

// HttpEntry todo
func (s *service) HTTPEntry() *http.EntrySet {
	return application.HttpEntry()
}

func init() {
	pkg.RegistryService("application", Service)
}

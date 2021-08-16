package grpc

import (
	"fmt"

	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/logger/zap"
	"github.com/infraboard/mcube/pb/http"
	clientv3 "go.etcd.io/etcd/client/v3"

	"github.com/infraboard/workflow/api/pkg"
	"github.com/infraboard/workflow/api/pkg/action"
	"github.com/infraboard/workflow/api/pkg/pipeline"
	"github.com/infraboard/workflow/conf"
)

var (
	// Service 服务实例
	Service = &impl{}
)

// NewEtcdRegister 初始化一个基于etcd的实例注册器
func NewPipelineService() (pipeline.ServiceServer, error) {
	etcdR := new(impl)

	return etcdR, nil
}

type impl struct {
	pipeline.UnimplementedServiceServer

	client *clientv3.Client
	log    logger.Logger
	action action.ServiceServer
}

func (s *impl) Config() error {
	s.client = conf.C().Etcd.GetClient()
	s.log = zap.L().Named("Pipeline")

	if pkg.Action == nil {
		return fmt.Errorf("dependence action service is nil")
	}
	s.action = pkg.Action
	return nil
}

// HttpEntry todo
func (s *impl) HTTPEntry() *http.EntrySet {
	return pipeline.HttpEntry()
}

func (e *impl) Debug(log logger.Logger) {
	e.log = log
}

func init() {
	pkg.RegistryService("pipeline", Service)
}

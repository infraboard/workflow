package impl

import (
	clientv3 "go.etcd.io/etcd/client/v3"

	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/workflow/api/pkg/task"
)

type impl struct {
	client *clientv3.Client
	logger.Logger
	task.UnimplementedServiceServer
}

// NewEtcdRegister 初始化一个基于etcd的实例注册器
func NewPipelineService() (task.ServiceServer, error) {
	etcdR := new(impl)
	return etcdR, nil
}

func (e *impl) Debug(log logger.Logger) {
	e.Logger = log
}

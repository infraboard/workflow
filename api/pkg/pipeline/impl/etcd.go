package impl

import (
	"fmt"
	"time"

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
func NewEtcdRegister(endpoints []string, username, password string) (task.ServiceServer, error) {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: time.Duration(5) * time.Second,
		Username:    username,
		Password:    password,
	})
	if err != nil {
		return nil, fmt.Errorf("connect etcd error, %s", err)
	}
	etcdR := new(impl)
	etcdR.client = client
	return etcdR, nil
}

func (e *impl) Debug(log logger.Logger) {
	e.Logger = log
}

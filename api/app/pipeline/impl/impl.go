package impl

import (
	"context"
	"fmt"
	"sync"

	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/logger/zap"
	"github.com/infraboard/mcube/pb/http"
	clientv3 "go.etcd.io/etcd/client/v3"

	"github.com/infraboard/workflow/api/app/action"
	"github.com/infraboard/workflow/api/app/pipeline"
	"github.com/infraboard/workflow/conf"
)

var (
	// Service 服务实例
	svr = &impl{
		watchCancel: make(map[int64]context.CancelFunc),
	}
)

type impl struct {
	pipeline.UnimplementedServiceServer

	client *clientv3.Client
	log    logger.Logger
	action action.ServiceServer

	watchCancel   map[int64]context.CancelFunc
	currentNumber int64
	l             sync.Mutex
}

func (i *impl) SetWatcherCancelFn(fn context.CancelFunc) int64 {
	i.l.Lock()
	defer i.l.Unlock()

	i.currentNumber++
	i.watchCancel[i.currentNumber] = fn
	return i.currentNumber
}

func (s *impl) Config() error {
	s.log = zap.L().Named("Pipeline")

	s.action = nil
	return nil
}

func (e *impl) Debug(log logger.Logger) {
	e.log = log
}

func (s *service) Name() string {
	return pipeline.AppName
}

func (s *service) Registry(server *grpc.Server) {
	pipeline.RegisterServiceServer(server, svr)
}

func init() {
	app.RegistryGrpcApp(svr)
}

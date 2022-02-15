package impl

import (
	"context"
	"sync"

	"github.com/infraboard/mcube/app"
	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/logger/zap"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"

	"github.com/infraboard/workflow/api/apps/action"
	"github.com/infraboard/workflow/api/apps/pipeline"
	"github.com/infraboard/workflow/conf"
)

var (
	// Service 服务实例
	svr = &impl{
		watchCancel: make(map[int64]context.CancelFunc),
	}
)

type impl struct {
	client *clientv3.Client
	log    logger.Logger
	action action.ServiceServer

	watchCancel   map[int64]context.CancelFunc
	currentNumber int64
	l             sync.Mutex

	pipeline.UnimplementedServiceServer
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
	s.client = conf.C().Etcd.GetClient()
	s.action = app.GetGrpcApp(action.AppName).(action.ServiceServer)
	return nil
}

func (e *impl) Debug(log logger.Logger) {
	e.log = log
}

func (s *impl) Name() string {
	return pipeline.AppName
}

func (s *impl) Registry(server *grpc.Server) {
	pipeline.RegisterServiceServer(server, svr)
}

func init() {
	app.RegistryGrpcApp(svr)
}

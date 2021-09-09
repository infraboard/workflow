package grpc

import (
	"context"
	"fmt"
	"sync"

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
	Service = &impl{
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

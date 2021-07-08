package etcd

import (
	"context"
	"errors"

	"github.com/infraboard/mcube/logger"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"

	"github.com/infraboard/workflow/api/pkg/pipeline"
	"github.com/infraboard/workflow/node/informer"
)

type shared struct {
	log           logger.Logger
	client        clientv3.Watcher
	stepHandler   informer.StepEventHandler
	stepWatchChan clientv3.WatchChan
}

func (i *shared) AddStepEventHandler(h informer.StepEventHandler) {
	i.stepHandler = h
}

// Run 启动 Watch
func (i *shared) Run(ctx context.Context) error {
	// 是否准备完成
	if err := i.isReady(); err != nil {
		return err
	}

	// 监听事件
	i.watchAll(ctx)

	go i.dealEvent()
	return nil
}

func (i *shared) dealEvent() {
	// 处理所有事件
	for {
		select {
		case stepResp := <-i.stepWatchChan:
			for _, event := range stepResp.Events {
				if err := i.notifyStep(event, stepResp.Header.GetRevision()); err != nil {
					i.log.Error(err)
				}
			}
		}
	}
}

func (i *shared) isReady() error {
	if i.stepHandler == nil {
		return errors.New("StepEventHandler not add")
	}
	return nil
}

func (i *shared) watchAll(ctx context.Context) {
	// 监听事件
	stepWatchKey := pipeline.EtcdStepPrefix()
	i.stepWatchChan = i.client.Watch(ctx, stepWatchKey, clientv3.WithPrefix())
	i.log.Infof("watch etcd step resource key: %s", stepWatchKey)
}

func (i *shared) notifyStep(event *clientv3.Event, eventVersion int64) error {
	// 解析对象
	obj, err := pipeline.LoadStepFromBytes(event.Kv.Value)
	if err != nil {
		return err
	}
	switch event.Type {
	case mvccpb.PUT:
		i.stepHandler.OnAdd(obj)
	case mvccpb.DELETE:
		i.stepHandler.OnDelete(obj)
	default:
	}
	return nil
}

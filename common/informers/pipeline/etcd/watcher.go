package etcd

import (
	"context"
	"errors"

	"github.com/infraboard/mcube/logger"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"

	"github.com/infraboard/workflow/api/pkg/pipeline"
	informer "github.com/infraboard/workflow/common/informers/pipeline"
)

type shared struct {
	log               logger.Logger
	client            clientv3.Watcher
	pipelineHandler   informer.PipelineEventHandler
	pipelineWatchChan clientv3.WatchChan
}

// AddPipelineEventHandler 添加事件处理回调
func (i *shared) AddPipelineTaskEventHandler(h informer.PipelineEventHandler) {
	i.pipelineHandler = h
}

// Run 启动 Watch
func (i *shared) Run(ctx context.Context) error {
	// 是否准备完成
	if err := i.isReady(); err != nil {
		return err
	}

	// 监听事件
	i.watch(ctx)

	// 后台处理事件
	go i.dealEvents()
	return nil
}

func (i *shared) dealEvents() {
	// 处理所有事件
	for {
		select {
		case ppResp := <-i.pipelineWatchChan:
			for _, event := range ppResp.Events {
				if err := i.notifyPipeline(event, ppResp.Header.GetRevision()); err != nil {
					i.log.Error(err)
				}
			}
		}
	}
}

func (i *shared) isReady() error {
	if i.pipelineHandler == nil {
		return errors.New("PipelineEventHandler not add")
	}
	return nil
}

func (i *shared) watch(ctx context.Context) {
	ppWatchKey := pipeline.EtcdPipelinePrefix()
	i.pipelineWatchChan = i.client.Watch(ctx, ppWatchKey, clientv3.WithPrefix())
	i.log.Infof("watch etcd pipeline resource key: %s", ppWatchKey)
}

func (i *shared) notifyPipeline(event *clientv3.Event, eventVersion int64) error {
	// 解析对象
	obj, err := pipeline.LoadPipelineFromBytes(event.Kv.Value)
	if err != nil {
		return err
	}
	obj.ResourceVersion = eventVersion
	switch event.Type {
	case mvccpb.PUT:
		i.pipelineHandler.OnAdd(obj)
	case mvccpb.DELETE:
		i.pipelineHandler.OnDelete(obj)
	default:
	}
	return nil
}

package etcd

import (
	"context"
	"errors"

	"github.com/infraboard/mcube/logger"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"

	"github.com/infraboard/workflow/api/pkg/pipeline"
	"github.com/infraboard/workflow/common/cache"
	informer "github.com/infraboard/workflow/common/informers/pipeline"
)

type shared struct {
	log               logger.Logger
	client            clientv3.Watcher
	indexer           cache.Indexer
	handler           informer.PipelineEventHandler
	filter            informer.PipelineFilterHandler
	pipelineWatchChan clientv3.WatchChan
}

// AddPipelineEventHandler 添加事件处理回调
func (i *shared) AddPipelineTaskEventHandler(h informer.PipelineEventHandler) {
	i.handler = h
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
	if i.handler == nil {
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
	i.log.Debugf("receive pipeline notify event, %s", event.Kv.Key)

	// 解析对象
	new, err := pipeline.LoadPipelineFromBytes(event.Kv.Value)
	if err != nil {
		return err
	}
	new.ResourceVersion = eventVersion

	old, hasOld, err := i.indexer.GetByKey(new.MakeObjectKey())
	if err != nil {
		return err
	}

	if i.filter != nil {
		if err := i.filter(new); err != nil {
			return err
		}
	}

	switch event.Type {
	case mvccpb.PUT:
		// 区分Update
		if hasOld {
			// 更新缓存
			i.log.Debugf("update pipeline: %s", new.ShortDescribe())
			if err := i.indexer.Update(new); err != nil {
				i.log.Errorf("update indexer cache error, %s", err)
			}
			i.handler.OnUpdate(old.(*pipeline.Pipeline), new)
		} else {
			// 添加缓存
			i.log.Debugf("add pipeline: %s", new.ShortDescribe())
			if err := i.indexer.Add(new); err != nil {
				i.log.Errorf("add indexer cache error, %s", err)
			}
			i.handler.OnAdd(new)
		}
	case mvccpb.DELETE:
		if !hasOld {
			return nil
		}
		// 清除缓存
		i.log.Debugf("delete pipeline: %s", new.ShortDescribe())
		if err := i.indexer.Delete(new); err != nil {
			i.log.Errorf("delete indexer cache error, %s", err)
		}
		i.handler.OnDelete(new)
	default:
	}
	return nil
}

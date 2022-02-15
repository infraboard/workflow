package etcd

import (
	"context"
	"errors"

	"github.com/infraboard/mcube/logger"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"

	"github.com/infraboard/workflow/api/apps/pipeline"
	"github.com/infraboard/workflow/common/cache"
	informer "github.com/infraboard/workflow/common/informers/step"
)

type shared struct {
	log       logger.Logger
	client    clientv3.Watcher
	indexer   cache.Indexer
	handler   informer.StepEventHandler
	filter    informer.StepFilterHandler
	watchChan clientv3.WatchChan
}

func (i *shared) AddStepEventHandler(h informer.StepEventHandler) {
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

	go i.dealEvent()
	return nil
}

func (i *shared) dealEvent() {
	// 处理所有事件
	for {
		select {
		case nodeResp := <-i.watchChan:
			for _, event := range nodeResp.Events {
				switch event.Type {
				case mvccpb.PUT:
					if err := i.handlePut(event, nodeResp.Header.GetRevision()); err != nil {
						i.log.Error(err)
					}
				case mvccpb.DELETE:
					if err := i.handleDelete(event); err != nil {
						i.log.Error(err)
					}
				default:
				}
			}
		}
	}
}

func (i *shared) isReady() error {
	if i.handler == nil {
		return errors.New("StepEventHandler not add")
	}
	return nil
}

func (i *shared) watch(ctx context.Context) {
	// 监听事件
	stepWatchKey := pipeline.EtcdStepPrefix()
	i.watchChan = i.client.Watch(ctx, stepWatchKey, clientv3.WithPrefix())
	i.log.Infof("watch etcd step resource key: %s", stepWatchKey)
}

func (i *shared) handlePut(event *clientv3.Event, eventVersion int64) error {
	i.log.Debugf("receive step put event, %s", event.Kv.Key)

	// 解析对象
	new, err := pipeline.LoadStepFromBytes(event.Kv.Value)
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

	// 区分Update
	if hasOld {
		// 更新缓存
		i.log.Debugf("update step store key: %s, status %s", new.Key, new.Status)
		if err := i.indexer.Update(new); err != nil {
			i.log.Errorf("update indexer cache error, %s", err)
		}
		i.handler.OnUpdate(old.(*pipeline.Step), new)
	} else {
		// 添加缓存
		i.log.Debugf("add step store key: %s, status %s", new.Key, new.Status)
		if err := i.indexer.Add(new); err != nil {
			i.log.Errorf("add indexer cache error, %s", err)
		}
		i.handler.OnAdd(new)
	}

	return nil
}

func (i *shared) handleDelete(event *clientv3.Event) error {
	key := event.Kv.Key
	i.log.Debugf("receive step delete event, %s", key)

	obj, ok, err := i.indexer.GetByKey(string(key))
	if err != nil {
		i.log.Errorf("get key %s from store error, %s", key)
	}
	if !ok {
		i.log.Warnf("key %s found in store", key)
	}
	if obj == nil {
		return nil
	}

	// 清除缓存
	if err := i.indexer.Delete(obj); err != nil {
		i.log.Errorf("delete indexer cache error, %s", err)
	}

	i.handler.OnDelete(obj.(*pipeline.Step))
	return nil
}

package etcd

import (
	"context"
	"errors"

	"github.com/infraboard/mcube/logger"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"

	"github.com/infraboard/workflow/api/pkg/pipeline"
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
	nodeName  string
}

func (i *shared) AddStepEventHandler(h informer.StepEventHandler) {
	i.handler = h
}

func (i *shared) AddStepFilterHandler(f informer.StepFilterHandler) {
	i.filter = f
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
		case stepResp := <-i.watchChan:
			for _, event := range stepResp.Events {
				if err := i.notifyStep(event, stepResp.Header.GetRevision()); err != nil {
					i.log.Error(err)
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

func (i *shared) notifyStep(event *clientv3.Event, eventVersion int64) error {
	// 解析对象
	new, err := pipeline.LoadStepFromBytes(event.Kv.Value)
	if err != nil {
		return err
	}
	old, hasOld, err := i.indexer.GetByKey(new.MakeObjectKey())
	if err != nil {
		return err
	}

	if i.filter != nil && !i.filter(new) {
		i.log.Debugf("step %s not match this node, now is %s, expect %s", new.String(), i.nodeName, new.ScheduledNodeName())
		return nil
	}

	switch event.Type {
	case mvccpb.PUT:
		// 区分Update
		if hasOld {
			// 更新缓存
			if err := i.indexer.Update(new); err != nil {
				i.log.Errorf("update indexer cache error, %s", err)
			}
			i.handler.OnUpdate(old.(*pipeline.Step), new)
		} else {
			// 添加缓存
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
		if err := i.indexer.Delete(new); err != nil {
			i.log.Errorf("delete indexer cache error, %s", err)
		}
		i.handler.OnDelete(new)
	default:
	}
	return nil
}

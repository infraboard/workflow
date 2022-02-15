package etcd

import (
	"context"
	"errors"
	"fmt"

	"github.com/infraboard/mcube/logger"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"

	"github.com/infraboard/workflow/api/apps/node"
	"github.com/infraboard/workflow/common/cache"
	informer "github.com/infraboard/workflow/common/informers/node"
)

type shared struct {
	log       logger.Logger
	client    clientv3.Watcher
	handler   informer.NodeEventHandler
	filter    informer.NodeFilterHandler
	watchChan clientv3.WatchChan
	indexer   cache.Indexer
}

// AddEventHandler 添加事件处理回调
func (i *shared) AddNodeEventHandler(h informer.NodeEventHandler) {
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
		return errors.New("NodeEventHandler not add")
	}
	return nil
}

func (s *shared) filt(node *node.Node) bool {
	if s.filter == nil {
		return false
	}

	err, ok := s.filter(node)
	if err == nil {
		s.log.Errorf("filt node error, %s", err)
		return false
	}

	return ok
}

func (i *shared) watch(ctx context.Context) {
	nodeWatchKey := node.EtcdNodePrefixWithType(node.NodeType)
	i.watchChan = i.client.Watch(ctx, nodeWatchKey, clientv3.WithPrefix())
	i.log.Infof("watch etcd node resource key: %s", nodeWatchKey)
}

func (i *shared) handlePut(event *clientv3.Event, eventVersion int64) error {
	i.log.Debugf("receive node put event, %s", event.Kv.Key)

	// 解析对象
	new, err := node.LoadNodeFromBytes(event.Kv.Value)
	if err != nil {
		return err
	}

	if new == nil {
		return fmt.Errorf("load node from bytes but get node object is nil")
	}

	old, hasOld, err := i.indexer.GetByKey(new.MakeObjectKey())
	if err != nil {
		return err
	}

	// 过滤掉不需要的
	if i.filt(new) {
		return nil
	}

	// 区分Update
	if hasOld {
		// 更新缓存
		i.log.Debugf("update node: %s", new.ShortDescribe())
		if err := i.indexer.Update(new); err != nil {
			i.log.Errorf("update indexer cache error, %s", err)
		}
		i.handler.OnUpdate(old.(*node.Node), new)
	} else {
		// 添加缓存
		i.log.Debugf("add node: %s", new.ShortDescribe())
		if err := i.indexer.Add(new); err != nil {
			i.log.Errorf("add indexer cache error, %s", err)
		}
		i.handler.OnAdd(new)
	}

	return nil
}

func (i *shared) handleDelete(event *clientv3.Event) error {
	key := event.Kv.Key
	i.log.Debugf("receive node delete event, %s", key)

	obj, ok, err := i.indexer.GetByKey(string(key))
	if err != nil {
		i.log.Errorf("get key %s from store error, %s", key)
	}
	if !ok {
		i.log.Warnf("key %s found in store", key)
	}

	// 清除缓存
	if err := i.indexer.Delete(obj); err != nil {
		i.log.Errorf("delete indexer cache error, %s", err)
	}

	i.handler.OnDelete(obj.(*node.Node))
	return nil
}

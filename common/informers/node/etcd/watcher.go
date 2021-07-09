package etcd

import (
	"context"
	"errors"
	"strings"

	"github.com/infraboard/mcube/logger"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"

	"github.com/infraboard/workflow/api/pkg/node"
	informer "github.com/infraboard/workflow/common/informers/node"
)

type shared struct {
	log           logger.Logger
	client        clientv3.Watcher
	nodeHandler   informer.NodeEventHandler
	nodeWatchChan clientv3.WatchChan
}

// AddEventHandler 添加事件处理回调
func (i *shared) AddNodeEventHandler(h informer.NodeEventHandler) {
	i.nodeHandler = h
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
		case nodeResp := <-i.nodeWatchChan:
			for _, event := range nodeResp.Events {
				if err := i.notifyNode(event, nodeResp.Header.GetRevision()); err != nil {
					i.log.Error(err)
				}
			}
		}
	}
}

func (i *shared) isReady() error {
	if i.nodeHandler == nil {
		return errors.New("NodeEventHandler not add")
	}
	return nil
}

func (i *shared) watch(ctx context.Context) {
	nodeWatchKey := node.EtcdNodePrefixWithType(node.NodeType)
	i.nodeWatchChan = i.client.Watch(ctx, nodeWatchKey, clientv3.WithPrefix())
	i.log.Infof("watch etcd node resource key: %s", nodeWatchKey)
}

func (i *shared) notifyNode(event *clientv3.Event, eventVersion int64) error {
	// 解析对象
	obj, err := node.LoadNodeFromBytes(event.Kv.Value)
	if err != nil {
		return err
	}

	// 解析事件为删除事件的时候，value内容为空，从key中解析 node 的 serviceName and instanceName
	if obj == nil {
		obj = new(node.Node)
		etcdNodeKeyList := strings.Split(string(event.Kv.Key), "/")
		obj.InstanceName = etcdNodeKeyList[len(etcdNodeKeyList)-1]
		obj.ServiceName = etcdNodeKeyList[1]
	}
	obj.ResourceVersion = eventVersion
	switch event.Type {
	case mvccpb.PUT:
		i.nodeHandler.OnAdd(obj)
	case mvccpb.DELETE:
		i.nodeHandler.OnDelete(obj)
	default:
	}
	return nil
}

package etcd

import (
	"context"

	"github.com/infraboard/mcube/logger"
	clientv3 "go.etcd.io/etcd/client/v3"

	"github.com/infraboard/workflow/api/pkg/node"
)

type lister struct {
	log    logger.Logger
	client clientv3.KV
}

// List 获取所有Node对象
func (l *lister) List(ctx context.Context) (ret []*node.Node, err error) {
	listKey := node.EtcdNodePrefix()
	l.log.Infof("list etcd node resource key: %s", listKey)
	resp, err := l.client.Get(ctx, listKey, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}
	nodes := []*node.Node{}
	for i := range resp.Kvs {
		// 解析对象
		node, err := node.LoadNodeFromBytes(resp.Kvs[i].Value)
		if err != nil {
			l.log.Error(err)
			continue
		}
		node.ResourceVersion = resp.Header.Revision
		nodes = append(nodes, node)
	}

	l.log.Infof("total nodes: %d", len(nodes))
	return nodes, nil
}

package etcd

import (
	"context"

	"github.com/infraboard/mcube/logger"
	clientv3 "go.etcd.io/etcd/client/v3"

	"github.com/infraboard/workflow/api/pkg/node"
	"github.com/infraboard/workflow/api/pkg/pipeline"
)

type lister struct {
	log    logger.Logger
	client clientv3.KV
}

func (l *lister) List(ctx context.Context) (ret []*pipeline.Step, err error) {
	listKey := node.EtcdNodePrefixWithType(node.NodeType)
	l.log.Infof("list etcd node resource key: %s", listKey)
	resp, err := l.client.Get(ctx, listKey, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}
	set := pipeline.NewStepSet()
	for i := range resp.Kvs {
		// 解析对象
		ins, err := pipeline.LoadStepFromBytes(resp.Kvs[i].Value)
		if err != nil {
			l.log.Error(err)
			continue
		}
		set.Add(ins)
	}

	return set.Items, nil
}
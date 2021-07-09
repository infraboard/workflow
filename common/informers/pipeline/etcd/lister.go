package etcd

import (
	"context"

	"github.com/infraboard/mcube/logger"
	clientv3 "go.etcd.io/etcd/client/v3"

	"github.com/infraboard/workflow/api/pkg/pipeline"
)

type lister struct {
	log    logger.Logger
	client clientv3.KV
}

func (l *lister) List(ctx context.Context, opts *pipeline.QueryPipelineOptions) (*pipeline.PipelineSet, error) {
	listKey := pipeline.EtcdPipelinePrefix()
	l.log.Infof("list etcd pipeline resource key: %s", listKey)
	resp, err := l.client.Get(ctx, listKey, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}
	ps := pipeline.NewPipelineSet()
	for i := range resp.Kvs {
		// 解析对象
		pt, err := pipeline.LoadPipelineFromBytes(resp.Kvs[i].Value)
		if err != nil {
			l.log.Errorf("load pipeline [key: %s, value: %s] error, %s", resp.Kvs[i].Key, string(resp.Kvs[i].Value), err)
			continue
		}

		pt.ResourceVersion = resp.Header.Revision
		ps.Add(pt)
	}
	return ps, nil
}

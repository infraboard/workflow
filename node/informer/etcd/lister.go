package etcd

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/infraboard/mcube/logger"
	clientv3 "go.etcd.io/etcd/client/v3"

	"github.com/infraboard/workflow/api/pkg/node"
	"github.com/infraboard/workflow/api/pkg/pipeline"
)

type lister struct {
	log    logger.Logger
	client clientv3.KV
}

// List 获取所有Node对象
func (l *lister) ListStep(ctx context.Context) (ret []*pipeline.Step, err error) {
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

func (l *lister) UpdateStep(step *pipeline.Step) error {
	objKey := pipeline.StepObjectKey(step.Key)
	objValue, err := json.Marshal(step)
	if err != nil {
		return err
	}
	if _, err := l.client.Put(context.Background(), objKey, string(objValue)); err != nil {
		return fmt.Errorf("update pipeline step '%s' to etcd3 failed: %s", objKey, err.Error())
	}
	return nil
}

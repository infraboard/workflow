package etcd

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/infraboard/mcube/logger"
	clientv3 "go.etcd.io/etcd/client/v3"

	"github.com/infraboard/workflow/api/pkg/node"
	"github.com/infraboard/workflow/api/pkg/pipeline"
	"github.com/infraboard/workflow/api/pkg/task"
	"github.com/infraboard/workflow/scheduler/informer"
)

type lister struct {
	log    logger.Logger
	client clientv3.KV
	prefix string
}

// List 获取所有Node对象
func (l *lister) ListNode(ctx context.Context) (ret []*node.Node, err error) {
	listKey := node.EtcdNodePrefixWithType(l.prefix, node.NodeType)
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
	return nodes, nil
}

func (l *lister) ListPipelineTask(ctx context.Context, opts *informer.QueryPipelineTaskOptions) (*task.PipelineTaskSet, error) {
	listKey := pipeline.EtcdPipelinePrefix(l.prefix)
	l.log.Infof("list etcd pipeline resource key: %s", listKey)
	resp, err := l.client.Get(ctx, listKey, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}
	ps := task.NewPipelineTaskSet()
	for i := range resp.Kvs {
		// 解析对象
		pt, err := task.LoadPipelineTaskFromBytes(resp.Kvs[i].Value)
		if err != nil {
			l.log.Error(err)
			continue
		}

		pt.ResourceVersion = resp.Header.Revision
		ps.Add(pt)
	}
	return ps, nil
}

func (l *lister) UpdateStep(step *pipeline.Step) error {
	objKey := step.EtcdObjectKey(pipeline.EtcdPipelinePrefix(l.prefix))
	objValue, err := json.Marshal(step)
	if err != nil {
		return err
	}
	if _, err := l.client.Put(context.Background(), objKey, string(objValue)); err != nil {
		return fmt.Errorf("update pipeline step '%s' to etcd3 failed: %s", objKey, err.Error())
	}
	return nil
}

func (l *lister) UpdateTask(t *task.PipelineTask) error {
	objKey := t.EtcdObjectKey(pipeline.EtcdPipelinePrefix(l.prefix))
	objValue, err := json.Marshal(t)
	if err != nil {
		return err
	}
	if _, err := l.client.Put(context.Background(), objKey, string(objValue)); err != nil {
		return fmt.Errorf("update pipeline task '%s' to etcd3 failed: %s", objKey, err.Error())
	}
	return nil
}

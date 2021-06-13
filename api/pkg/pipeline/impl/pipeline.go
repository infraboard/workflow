package impl

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/infraboard/workflow/api/pkg/pipeline"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func (i *impl) CreatePipeline(ctx context.Context, req *pipeline.CreatePipelineRequest) (
	*pipeline.Pipeline, error) {

	p, err := pipeline.NewPipeline(req)
	if err != nil {
		return nil, err
	}

	value, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}

	objKey := p.EtcdObjectKey(i.prefix)
	objValue := string(value)

	if _, err := i.client.Put(context.Background(), objKey, objValue); err != nil {
		return nil, fmt.Errorf("put pipeline with key: %s, error, %s", objKey, err.Error())
	}
	i.log.Debugf("create pipeline success, key: %s", objKey)
	return p, nil
}

func (i *impl) QueryPipeline(ctx context.Context, req *pipeline.QueryPipelineRequest) (
	*pipeline.PipelineSet, error) {
	listKey := pipeline.EtcdPipelinePrefix(i.prefix)
	i.log.Infof("list etcd pipeline resource key: %s", listKey)
	resp, err := i.client.Get(ctx, listKey, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}

	ps := pipeline.NewPipelineSet()
	for index := range resp.Kvs {
		// 解析对象
		ins, err := pipeline.LoadPipelineFromBytes(resp.Kvs[index].Value)
		if err != nil {
			i.log.Error(err)
			continue
		}
		ins.ResourceVersion = resp.Header.Revision
		ps.Add(ins)
	}
	return ps, nil
}

func (i *impl) CreateAction(context.Context, *pipeline.CreateActionRequest) (
	*pipeline.Action, error) {
	return nil, nil
}

func (i *impl) QueryAction(context.Context, *pipeline.QueryActionRequest) (
	*pipeline.ActionSet, error) {
	return nil, nil
}

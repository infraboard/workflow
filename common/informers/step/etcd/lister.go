package etcd

import (
	"context"

	"github.com/infraboard/mcube/exception"
	"github.com/infraboard/mcube/logger"
	clientv3 "go.etcd.io/etcd/client/v3"

	"github.com/infraboard/workflow/api/pkg/pipeline"
	"github.com/infraboard/workflow/common/informers/step"
)

type lister struct {
	log    logger.Logger
	client clientv3.KV
	filter step.StepFilterHandler
}

func (l *lister) List(ctx context.Context) (ret []*pipeline.Step, err error) {
	listKey := pipeline.EtcdStepPrefix()

	l.log.Infof("list etcd step resource key: %s", listKey)
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

		if l.filter != nil {
			if err := l.filter(ins); err != nil {
				l.log.Error(err)
				continue
			}
		}

		ins.ResourceVersion = resp.Header.Revision
		set.Add(ins)
	}

	return set.Items, nil
}

func (l *lister) Get(ctx context.Context, key string) (*pipeline.Step, error) {
	descKey := pipeline.StepObjectKey(key)
	l.log.Infof("describe etcd step resource key: %s", descKey)
	resp, err := l.client.Get(ctx, descKey)
	if err != nil {
		return nil, err
	}

	if resp.Count == 0 {
		return nil, nil
	}

	if resp.Count > 1 {
		return nil, exception.NewInternalServerError("step find more than one: %d", resp.Count)
	}

	ins := pipeline.NewDefaultStep()
	for index := range resp.Kvs {
		// 解析对象
		ins, err = pipeline.LoadStepFromBytes(resp.Kvs[index].Value)
		if err != nil {
			l.log.Error(err)
			continue
		}
	}
	return ins, nil
}

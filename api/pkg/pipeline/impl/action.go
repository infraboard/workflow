package impl

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/infraboard/keyauth/client/session"
	"github.com/infraboard/mcube/exception"
	"github.com/infraboard/mcube/grpc/gcontext"
	"github.com/infraboard/workflow/api/pkg/pipeline"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func (i *impl) CreateAction(ctx context.Context, req *pipeline.CreateActionRequest) (
	*pipeline.Action, error) {

	in, err := gcontext.GetGrpcInCtx(ctx)
	if err != nil {
		return nil, err
	}
	tk := session.S().GetToken(in.GetAccessToKen())
	if tk == nil {
		return nil, exception.NewUnauthorized("token required")
	}

	p, err := pipeline.NewAction(req)
	if err != nil {
		return nil, err
	}

	p.NeedSecret = req.NeedSecret
	p.CreateBy = tk.Account
	p.Domain = tk.Domain
	p.Namespace = tk.Namespace

	value, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}

	objKey := p.EtcdObjectKey()
	objValue := string(value)

	if _, err := i.client.Put(context.Background(), objKey, objValue); err != nil {
		return nil, fmt.Errorf("put action with key: %s, error, %s", objKey, err.Error())
	}
	i.log.Debugf("create action success, key: %s", objKey)

	return p, nil
}

func (i *impl) QueryAction(ctx context.Context, req *pipeline.QueryActionRequest) (
	*pipeline.ActionSet, error) {

	listKey := pipeline.EtcdActionPrefix()
	i.log.Debugf("list etcd action resource key: %s", listKey)
	resp, err := i.client.Get(ctx, listKey, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}

	ps := pipeline.NewActionSet()
	for index := range resp.Kvs {
		// 解析对象
		ins, err := pipeline.LoadActionFromBytes(resp.Kvs[index].Value)
		if err != nil {
			i.log.Error(err)
			continue
		}
		ps.Add(ins)
	}
	return ps, nil
}

func (i *impl) DeleteAction(ctx context.Context, req *pipeline.DeleteActionRequest) (
	*pipeline.Action, error) {

	i.log.Debug("req:", req)
	ac := &pipeline.Action{}
	objKey := pipeline.EtcdActionPrefix()
	i.log.Infof("list etcd action resource key: %s", objKey)
	resp, err := i.client.Delete(ctx, objKey, clientv3.WithPrefix())
	if err != nil {
		return ac, fmt.Errorf("put action with key: %s, error, %s", objKey, err.Error())
	}

	i.log.Debugf("create action success, key num: %D", resp.Deleted)
	return ac, nil
}

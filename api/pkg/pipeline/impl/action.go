package impl

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/infraboard/keyauth/client/session"
	"github.com/infraboard/mcube/exception"
	"github.com/infraboard/mcube/grpc/gcontext"
	"github.com/infraboard/mcube/pb/resource"
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

	a, err := pipeline.NewAction(req)
	if err != nil {
		return nil, err
	}

	a.UpdateOwner(tk)

	value, err := json.Marshal(a)
	if err != nil {
		return nil, err
	}

	objKey := a.EtcdObjectKey()
	objValue := string(value)

	if _, err := i.client.Put(context.Background(), objKey, objValue); err != nil {
		return nil, fmt.Errorf("put action with key: %s, error, %s", objKey, err.Error())
	}
	i.log.Debugf("create action success, key: %s", objKey)

	return a, nil
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

func (i *impl) DescribeAction(ctx context.Context, req *pipeline.DescribeActionRequest) (
	*pipeline.Action, error) {
	in, err := gcontext.GetGrpcInCtx(ctx)
	if err != nil {
		return nil, err
	}
	tk := session.S().GetToken(in.GetAccessToKen())
	if tk == nil {
		return nil, exception.NewUnauthorized("token required")
	}

	var ins *pipeline.Action
	// 先搜索Namespace内
	ins, err = i.describeAction(ctx, tk.Namespace, req)
	if err != nil {
		return nil, err
	}

	// 再搜索全局
	if ins == nil {
		ins, err = i.describeAction(ctx, resource.VisiableMode_GLOBAL.String(), req)
		if err != nil {
			return nil, err
		}
	}

	if ins == nil {
		return nil, exception.NewNotFound("action %s not found", req.Name)
	}

	return ins, nil
}

func (i *impl) describeAction(ctx context.Context, namespace string, req *pipeline.DescribeActionRequest) (
	*pipeline.Action, error) {
	descKey := pipeline.ActionObjectKey(namespace, req.Name)
	i.log.Infof("describe etcd action resource key: %s", descKey)
	resp, err := i.client.Get(ctx, descKey)
	if err != nil {
		return nil, err
	}

	if resp.Count == 0 {
		return nil, nil
	}

	if resp.Count > 1 {
		return nil, exception.NewInternalServerError("action find more than one: %d", resp.Count)
	}

	ins := pipeline.NewDefaultAction()
	for index := range resp.Kvs {
		// 解析对象
		ins, err = pipeline.LoadActionFromBytes(resp.Kvs[index].Value)
		if err != nil {
			i.log.Error(err)
			continue
		}
	}
	return ins, nil
}

func (i *impl) DeleteAction(ctx context.Context, req *pipeline.DeleteActionRequest) (
	*pipeline.Action, error) {
	in, err := gcontext.GetGrpcInCtx(ctx)
	if err != nil {
		return nil, err
	}
	tk := session.S().GetToken(in.GetAccessToKen())
	if tk == nil {
		return nil, exception.NewUnauthorized("token required")
	}
	descKey := pipeline.ActionObjectKey(req.Namespace(tk), req.Name)
	i.log.Infof("delete etcd action resource key: %s", descKey)
	resp, err := i.client.Delete(ctx, descKey, clientv3.WithPrevKV())
	if err != nil {
		return nil, err
	}

	if resp.Deleted == 0 {
		return nil, exception.NewNotFound("action %s not found", req.Name)
	}

	ins := pipeline.NewDefaultAction()
	for index := range resp.PrevKvs {
		// 解析对象
		ins, err = pipeline.LoadActionFromBytes(resp.PrevKvs[index].Value)
		if err != nil {
			i.log.Error(err)
			continue
		}
		ins.ResourceVersion = resp.Header.Revision
	}
	return ins, nil
}

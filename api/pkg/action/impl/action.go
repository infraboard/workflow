package impl

import (
	"context"

	"github.com/infraboard/keyauth/client/session"
	"github.com/infraboard/mcube/exception"
	"github.com/infraboard/mcube/grpc/gcontext"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/infraboard/workflow/api/pkg/action"
)

func (i *impl) CreateAction(ctx context.Context, req *action.CreateActionRequest) (
	*action.Action, error) {
	in, err := gcontext.GetGrpcInCtx(ctx)
	if err != nil {
		return nil, err
	}
	tk := session.S().GetToken(in.GetAccessToKen())
	if tk == nil {
		return nil, exception.NewUnauthorized("token required")
	}

	a, err := action.NewAction(req)
	if err != nil {
		return nil, err
	}

	a.UpdateOwner(tk)

	if _, err := i.col.InsertOne(context.TODO(), a); err != nil {
		return nil, exception.NewInternalServerError("inserted a action document error, %s", err)
	}

	return a, nil
}

func (i *impl) QueryAction(ctx context.Context, req *action.QueryActionRequest) (
	*action.ActionSet, error) {

	query := newQueryActionRequest(req)
	resp, err := i.col.Find(context.TODO(), query.FindFilter(), query.FindOptions())

	if err != nil {
		return nil, exception.NewInternalServerError("find action error, error is %s", err)
	}

	set := action.NewActionSet()
	// 循环
	for resp.Next(context.TODO()) {
		a := action.NewDefaultAction()
		if err := resp.Decode(a); err != nil {
			return nil, exception.NewInternalServerError("decode action error, error is %s", err)
		}

		set.Add(a)
	}

	// count
	count, err := i.col.CountDocuments(context.TODO(), query.FindFilter())
	if err != nil {
		return nil, exception.NewInternalServerError("get action count error, error is %s", err)
	}
	set.Total = count

	return set, nil
}

func (i *impl) DescribeAction(ctx context.Context, req *action.DescribeActionRequest) (
	*action.Action, error) {
	if req.Namespace == "" {
		in, err := gcontext.GetGrpcInCtx(ctx)
		if err != nil {
			return nil, err
		}

		tk := session.S().GetToken(in.GetAccessToKen())
		if tk == nil {
			return nil, exception.NewUnauthorized("token required")
		}
		req.Namespace = tk.Namespace
	}

	if err := req.Validate(); err != nil {
		return nil, exception.NewBadRequest("validate DescribeActionRequest error, %s", err)
	}

	desc, err := newDescActionRequest(req)
	if err != nil {
		return nil, exception.NewBadRequest(err.Error())
	}

	ins := action.NewDefaultAction()
	if err := i.col.FindOne(context.TODO(), desc.FindFilter()).Decode(ins); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, exception.NewNotFound("action %s not found", req)
		}

		return nil, exception.NewInternalServerError("find domain %s error, %s", req.Name, err)
	}

	return ins, nil
}

func (i *impl) UpdateAction(context.Context, *action.UpdateActionRequest) (*action.Action, error) {
	return nil, nil
}

func (i *impl) DeleteAction(ctx context.Context, req *action.DeleteActionRequest) (
	*action.Action, error) {
	in, err := gcontext.GetGrpcInCtx(ctx)
	if err != nil {
		return nil, err
	}
	tk := session.S().GetToken(in.GetAccessToKen())
	if tk == nil {
		return nil, exception.NewUnauthorized("token required")
	}

	// descKey := pipeline.ActionObjectKey(req.Namespace(tk), req.Name, req.Version)
	// i.log.Infof("delete etcd action resource key: %s", descKey)
	// resp, err := i.client.Delete(ctx, descKey, clientv3.WithPrevKV())
	// if err != nil {
	// 	return nil, err
	// }

	// if resp.Deleted == 0 {
	// 	return nil, exception.NewNotFound("action %s not found", req.Name)
	// }

	// ins := pipeline.NewDefaultAction()
	// for index := range resp.PrevKvs {
	// 	// 解析对象
	// 	ins, err = pipeline.LoadActionFromBytes(resp.PrevKvs[index].Value)
	// 	if err != nil {
	// 		i.log.Error(err)
	// 		continue
	// 	}
	// 	ins.ResourceVersion = resp.Header.Revision
	// }
	return nil, nil
}

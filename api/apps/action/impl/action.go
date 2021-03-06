package impl

import (
	"context"
	"time"

	"github.com/infraboard/mcube/exception"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/infraboard/workflow/api/apps/action"
)

func (i *service) CreateAction(ctx context.Context, req *action.CreateActionRequest) (
	*action.Action, error) {
	a, err := action.NewAction(req)
	if err != nil {
		return nil, err
	}

	// 获取之前最新的版本
	if _, err := i.col.InsertOne(context.TODO(), a); err != nil {
		return nil, exception.NewInternalServerError("inserted a action document error, %s", err)
	}

	return a, nil
}

func (i *service) QueryAction(ctx context.Context, req *action.QueryActionRequest) (
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

func (i *service) DescribeAction(ctx context.Context, req *action.DescribeActionRequest) (
	*action.Action, error) {
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

		return nil, exception.NewInternalServerError("find action %s error, %s", req.Name, err)
	}

	return ins, nil
}

func (i *service) UpdateAction(ctx context.Context, req *action.UpdateActionRequest) (*action.Action, error) {
	if err := req.Validate(); err != nil {
		return nil, exception.NewBadRequest(err.Error())
	}

	ins, err := i.DescribeAction(ctx, action.NewDescribeActionRequest(req.Name, req.Version))
	if err != nil {
		return nil, err
	}

	ins.Update(req)
	ins.UpdateAt = time.Now().UnixMilli()
	_, err = i.col.UpdateOne(context.TODO(), bson.M{"name": req.Name, "version": req.Version}, bson.M{"$set": ins})
	if err != nil {
		return nil, exception.NewInternalServerError("update action(%s) error, %s", ins.Key(), err)
	}

	return ins, nil
}

func (i *service) DeleteAction(ctx context.Context, req *action.DeleteActionRequest) (
	*action.Action, error) {
	ins, err := i.DescribeAction(ctx, action.NewDescribeActionRequest(req.Name, req.Version))
	if err != nil {
		return nil, err
	}

	if ins.Namespace != req.Namespace {
		return nil, exception.NewBadRequest("target action namespace not match your namespace")
	}

	delReq, err := newDeleteActionRequest(req)
	if err != nil {
		return nil, exception.NewBadRequest(err.Error())
	}

	if _, err := i.col.DeleteOne(context.TODO(), delReq.DeleteFilter()); err != nil {
		return nil, err
	}

	return ins, nil
}

package grpc

import (
	"context"

	"github.com/infraboard/mcube/exception"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/infraboard/workflow/api/pkg/application"
	"github.com/infraboard/workflow/api/pkg/deploy"
)

func (s *service) CreateApplicationDeploy(ctx context.Context, req *deploy.CreateApplicationDeployRequest) (
	*deploy.ApplicationDeploy, error) {
	ins, err := deploy.NewApplicationDeploy(req)
	if err != nil {
		return nil, exception.NewBadRequest("validate CreateApplicationDeployRequest error, %s", err)
	}

	_, err = s.app.DescribeApplication(ctx, application.NewDescribeApplicationRequestWithID(req.AppId))
	if err != nil {
		return nil, err
	}

	if _, err := s.col.InsertOne(context.TODO(), ins); err != nil {
		return nil, exception.NewInternalServerError("inserted a deploy document error, %s", err)
	}

	return ins, nil
}

func (s *service) QueryApplicationDeploy(ctx context.Context, req *deploy.QueryApplicationDeployRequest) (
	*deploy.ApplicationDeploySet, error) {
	query := newQueryApplicationDeployRequest(req)
	resp, err := s.col.Find(context.TODO(), query.FindFilter(), query.FindOptions())

	if err != nil {
		return nil, exception.NewInternalServerError("find deploy error, error is %s", err)
	}

	set := deploy.NewApplicationDeploySet()
	// 循环
	for resp.Next(context.TODO()) {
		ins := deploy.NewDefaultApplicationDeploy()
		if err := resp.Decode(ins); err != nil {
			return nil, exception.NewInternalServerError("decode deploy error, error is %s", err)
		}

		ins.Desensitize()
		set.Add(ins)
	}

	// count
	count, err := s.col.CountDocuments(context.TODO(), query.FindFilter())
	if err != nil {
		return nil, exception.NewInternalServerError("get deploy count error, error is %s", err)
	}
	set.Total = count
	return set, nil
}

func (s *service) DescribeApplicationDeploy(ctx context.Context, req *deploy.DescribeApplicationDeployRequest) (
	*deploy.ApplicationDeploy, error) {
	if err := req.Validate(); err != nil {
		return nil, exception.NewBadRequest("validate DescribeApplicationDeployRequest error, %s", err)
	}

	desc, err := newDescRequest(req)
	if err != nil {
		return nil, exception.NewBadRequest(err.Error())
	}

	ins := deploy.NewDefaultApplicationDeploy()
	if err := s.col.FindOne(context.TODO(), desc.FindFilter()).Decode(ins); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, exception.NewNotFound("deploy %s not found", req)
		}

		return nil, exception.NewInternalServerError("find deploy %s error, %s", req.Id, err)
	}
	return ins, nil
}

func (s *service) DeleteApplicationDeploy(ctx context.Context, req *deploy.DeleteApplicationDeployRequest) (
	*deploy.ApplicationDeploy, error) {
	ins, err := s.DescribeApplicationDeploy(ctx, deploy.NewDescribeApplicationDeployRequestWithID(req.Id))
	if err != nil {
		return nil, err
	}

	if _, err := s.col.DeleteOne(context.TODO(), bson.M{"_id": ins.Id}); err != nil {
		return nil, err
	}

	return ins, nil
}

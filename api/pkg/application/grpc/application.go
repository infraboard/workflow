package grpc

import (
	"context"
	"fmt"

	"github.com/infraboard/mcube/exception"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/infraboard/workflow/api/pkg/application"
	"github.com/infraboard/workflow/common/repo/gitlab"
)

func (s *service) CreateApplication(ctx context.Context, req *application.CreateApplicationRequest) (
	*application.Application, error) {
	ins, err := application.NewApplication(req)
	if err != nil {
		return nil, err
	}

	hookId, err := s.setWebHook(req, ins.GenWebHook(s.platform))
	if err != nil {
		ins.AddError(fmt.Errorf("add web hook error, %s", err))
	}
	ins.ScmHookId = hookId

	if _, err := s.col.InsertOne(context.TODO(), ins); err != nil {
		return nil, exception.NewInternalServerError("inserted a application document error, %s", err)
	}

	return ins, nil
}

func (s *service) setWebHook(req *application.CreateApplicationRequest, hook *gitlab.WebHook) (string, error) {
	if req.NeedSetHook() {
		return "", nil
	}

	addr, err := req.GetScmAddr()
	if err != nil {
		return "", fmt.Errorf("get scm addr from http_url error, %s", err)
	}

	repo := gitlab.NewRepository(addr, req.ScmPrivateToken)
	addHookReq := gitlab.NewAddProjectHookRequest(req.Int64ScmProjectID(), hook)
	resp, err := repo.AddProjectHook(addHookReq)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%d", resp.ID), nil
}

func (s *service) QueryApplication(ctx context.Context, req *application.QueryApplicationRequest) (
	*application.ApplicationSet, error) {
	query := newQueryApplicationRequest(req)
	resp, err := s.col.Find(context.TODO(), query.FindFilter(), query.FindOptions())

	if err != nil {
		return nil, exception.NewInternalServerError("find application error, error is %s", err)
	}

	set := application.NewApplicationSet()
	// 循环
	for resp.Next(context.TODO()) {
		a := application.NewDefaultApplication()
		if err := resp.Decode(a); err != nil {
			return nil, exception.NewInternalServerError("decode application error, error is %s", err)
		}

		set.Add(a)
	}

	// count
	count, err := s.col.CountDocuments(context.TODO(), query.FindFilter())
	if err != nil {
		return nil, exception.NewInternalServerError("get application count error, error is %s", err)
	}
	set.Total = count
	return set, nil
}

func (s *service) DescribeApplication(ctx context.Context, req *application.DescribeApplicationRequest) (
	*application.Application, error) {
	if err := req.Validate(); err != nil {
		return nil, exception.NewBadRequest("validate DescribeApplicationRequest error, %s", err)
	}

	desc, err := newDescRequest(req)
	if err != nil {
		return nil, exception.NewBadRequest(err.Error())
	}

	ins := application.NewDefaultApplication()
	if err := s.col.FindOne(context.TODO(), desc.FindFilter()).Decode(ins); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, exception.NewNotFound("application %s not found", req)
		}

		return nil, exception.NewInternalServerError("find application %s error, %s", req.Id, err)
	}
	return ins, nil
}

func (s *service) DeleteApplication(ctx context.Context, req *application.DeleteApplicationRequest) (
	*application.Application, error) {
	ins, err := s.DescribeApplication(ctx, application.NewDescribeApplicationRequestWithName(req.Namespace, req.Name))
	if err != nil {
		return nil, err
	}

	if _, err := s.col.DeleteOne(context.TODO(), bson.M{"_id": ins.Id}); err != nil {
		return nil, err
	}
	return ins, nil
}

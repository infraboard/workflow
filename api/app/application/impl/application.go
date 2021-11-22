package impl

import (
	"context"
	"fmt"

	"github.com/infraboard/mcube/exception"
	"github.com/infraboard/mcube/pb/request"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/infraboard/workflow/api/app/application"
	"github.com/infraboard/workflow/api/app/scm/gitlab"
)

func (s *service) CreateApplication(ctx context.Context, req *application.CreateApplicationRequest) (
	*application.Application, error) {
	ins, err := application.NewApplication(req)
	if err != nil {
		return nil, err
	}

	hookId, err := s.setWebHook(req, ins.GenWebHook(s.platform))
	if err != nil {
		ins.HookError = fmt.Sprintf("add web hook error, %s", err)
	}
	ins.ScmHookId = hookId

	if _, err := s.col.InsertOne(context.TODO(), ins); err != nil {
		return nil, exception.NewInternalServerError("inserted a application document error, %s", err)
	}

	return ins, nil
}

func (s *service) setWebHook(req *application.CreateApplicationRequest, hook *gitlab.WebHook) (string, error) {
	if !req.NeedSetHook() {
		return "", nil
	}

	addr, err := req.GetScmAddr()
	if err != nil {
		return "", fmt.Errorf("get scm addr from http_url error, %s", err)
	}

	repo := gitlab.NewSCM(addr, req.ScmPrivateToken)
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

		a.Desensitize()
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

	// 删除Hook
	if err := s.delWebHook(ins); err != nil {
		s.log.Errorf("delete scm hook error, %s", err)
	}

	if _, err := s.col.DeleteOne(context.TODO(), bson.M{"_id": ins.Id}); err != nil {
		return nil, err
	}
	return ins, nil
}

func (s *service) delWebHook(req *application.Application) error {
	if req.ScmHookId == "" {
		return nil
	}

	if req.ScmPrivateToken == "" {
		s.log.Errorf("delete scm hook error, scm_private_token is empty")
		return nil
	}

	addr, err := req.GetScmAddr()
	if err != nil {
		return fmt.Errorf("get scm addr from http_url error, %s", err)
	}

	repo := gitlab.NewSCM(addr, req.ScmPrivateToken)
	delHookReq := gitlab.NewDeleteProjectReqeust(req.Int64ScmProjectID(), req.Int64ScmHookID())

	return repo.DeleteProjectHook(delHookReq)
}

func (s *service) UpdateApplication(ctx context.Context, req *application.UpdateApplicationRequest) (
	*application.Application, error) {
	if err := req.Validate(); err != nil {
		return nil, exception.NewBadRequest("validate update application error, %s", err)
	}

	app, err := s.DescribeApplication(ctx, application.NewDescribeApplicationRequestWithID(req.Id))
	if err != nil {
		return nil, err
	}

	switch req.UpdateMode {
	case request.UpdateMode_PUT:
		app.Update(req.UpdateBy, req.Data)
	case request.UpdateMode_PATCH:
		app.Patch(req.UpdateBy, req.Data)
	default:
		return nil, fmt.Errorf("unknown update mode: %s", req.UpdateMode)
	}

	if err := s.update(ctx, app); err != nil {
		return nil, err
	}

	return app, nil
}

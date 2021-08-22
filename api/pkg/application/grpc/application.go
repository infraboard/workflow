package grpc

import (
	"context"
	"fmt"

	"github.com/infraboard/mcube/exception"

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

	return string(resp.ID), nil
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

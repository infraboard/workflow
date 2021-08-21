package grpc

import (
	"context"

	"github.com/infraboard/mcube/exception"

	"github.com/infraboard/workflow/api/pkg/application"
)

func (s *service) CreateApplication(ctx context.Context, req *application.CreateApplicationRequest) (
	*application.Application, error) {
	ins, err := application.NewApplication(req)
	if err != nil {
		return nil, err
	}

	if _, err := s.col.InsertOne(context.TODO(), ins); err != nil {
		return nil, exception.NewInternalServerError("inserted a application document error, %s", err)
	}

	return ins, nil
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

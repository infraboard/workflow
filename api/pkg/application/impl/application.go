package impl

import (
	"context"

	"github.com/infraboard/mcube/grpc/gcontext"

	"github.com/infraboard/workflow/api/pkg"
	"github.com/infraboard/workflow/api/pkg/application"
)

func (s *service) CreateApplication(ctx context.Context, req *application.CreateApplicationRequest) (
	*application.Application, error) {
	in, err := gcontext.GetGrpcInCtx(ctx)
	if err != nil {
		return nil, err
	}
	tk := pkg.S().GetToken(in.GetRequestID())
	s.log.Debug(tk)
	return application.NewApplication(req), nil
}

func (s *service) QueryApplication(ctx context.Context, req *application.QueryApplicationRequest) (
	*application.ApplicationSet, error) {
	return application.NewApplicationSet(), nil
}

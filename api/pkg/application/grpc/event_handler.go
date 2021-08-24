package grpc

import (
	"context"

	"github.com/infraboard/mcube/exception"

	"github.com/infraboard/workflow/api/pkg/application"
	"github.com/infraboard/workflow/api/pkg/pipeline"
)

func (s *service) HandleApplicationEvent(ctx context.Context, in *application.ApplicationEvent) (*pipeline.PipelineSet, error) {
	if err := in.Validate(); err != nil {
		return nil, exception.NewBadRequest("validate ApplicationEvent error, %s", err)
	}

	// 查询应用
	app, err := s.DescribeApplication(ctx, application.NewDescribeApplicationRequestWithID(in.AppId))
	if err != nil {
		return nil, err
	}

	// 找出匹配的pipline
	matched := app.MatchPipeline(in.WebhookEvent)
	if len(matched) == 0 {
		s.log.Infof("application %s no pipeline matched the event", app.Id, in.WebhookEvent.ShortDesc())
		return nil, nil
	}

	// pipeline参数实例化

	return nil, nil
}

package grpc

import (
	"context"
	"fmt"

	"github.com/infraboard/mcube/exception"

	"github.com/infraboard/workflow/api/pkg/application"
)

func (s *service) HandleApplicationEvent(ctx context.Context, in *application.ApplicationEvent) (
	*application.Application, error) {
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
		s.log.Infof("application %s no pipeline matched the event %s", app.Id, in.WebhookEvent.ShortDesc())
		return app, nil
	}

	// 运行这些匹配到的pipeline
	for i := range matched {
		req := matched[i]
		req.HookEvent = in.WebhookEvent
		req.Domain = app.Domain
		req.Namespace = app.Namespace
		req.CreateBy = fmt.Sprintf("@app:%s", app.Name)
		status := application.NewPipelineCreateStatus(req.Name)

		s.log.Debugf("start create pipeline: %s", req.Name)
		p, err := s.pipeline.CreatePipeline(ctx, req)
		if err != nil {
			status.CreateError = err.Error()
		} else {
			status.Pipeline = p
		}

		app.UpdatePipelineStatus(status)
	}

	// 更新应用状态
	if err := s.update(ctx, app); err != nil {
		return nil, err
	}

	return app, nil
}

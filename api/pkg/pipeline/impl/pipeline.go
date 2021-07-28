package impl

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/infraboard/keyauth/client/session"
	"github.com/infraboard/mcube/exception"
	"github.com/infraboard/mcube/grpc/gcontext"
	"github.com/infraboard/workflow/api/pkg/pipeline"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func (i *impl) CreatePipeline(ctx context.Context, req *pipeline.CreatePipelineRequest) (
	*pipeline.Pipeline, error) {
	in, err := gcontext.GetGrpcInCtx(ctx)
	if err != nil {
		return nil, err
	}
	tk := session.S().GetToken(in.GetAccessToKen())
	if tk == nil {
		return nil, exception.NewUnauthorized("token required")
	}

	p, err := pipeline.NewPipeline(req)
	if err != nil {
		return nil, err
	}
	p.UpdateOwner(tk)

	if err := i.validatePipeline(ctx, p); err != nil {
		return nil, err
	}

	value, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}

	objKey := p.EtcdObjectKey()
	objValue := string(value)

	if _, err := i.client.Put(context.Background(), objKey, objValue); err != nil {
		return nil, fmt.Errorf("put pipeline with key: %s, error, %s", objKey, err.Error())
	}
	i.log.Debugf("create pipeline success, key: %s", objKey)
	return p, nil
}

func (i *impl) validatePipeline(ctx context.Context, p *pipeline.Pipeline) error {
	for index := range p.Stages {
		stage := p.Stages[index]
		if err := i.validateStage(ctx, stage); err != nil {
			return err
		}
	}

	return nil
}

func (i *impl) validateStage(ctx context.Context, s *pipeline.Stage) error {
	if s.StepCount() == 0 {
		return fmt.Errorf("stage %s host no steps", s.ShortDesc())
	}

	for index := range s.Steps {
		step := s.Steps[index]
		if err := i.validateStep(ctx, step); err != nil {
			return err
		}
	}

	return nil
}

func (i *impl) validateStep(ctx context.Context, s *pipeline.Step) error {
	_, err := i.DescribeAction(ctx, pipeline.NewDescribeActionRequestWithName(s.Action))
	if err != nil {
		return err
	}
	return nil
}

func (i *impl) QueryPipeline(ctx context.Context, req *pipeline.QueryPipelineRequest) (
	*pipeline.PipelineSet, error) {
	listKey := pipeline.EtcdPipelinePrefix()
	i.log.Infof("list etcd pipeline resource key: %s", listKey)
	resp, err := i.client.Get(ctx, listKey, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}

	ps := pipeline.NewPipelineSet()
	for index := range resp.Kvs {
		// 解析对象
		ins, err := pipeline.LoadPipelineFromBytes(resp.Kvs[index].Value)
		if err != nil {
			i.log.Error(err)
			continue
		}
		ins.ResourceVersion = resp.Header.Revision
		ps.Add(ins)
	}
	return ps, nil
}

func (i *impl) DescribePipeline(ctx context.Context, req *pipeline.DescribePipelineRequest) (
	*pipeline.Pipeline, error) {
	in, err := gcontext.GetGrpcInCtx(ctx)
	if err != nil {
		return nil, err
	}

	if req.Namespace == "" {
		tk := session.S().GetToken(in.GetAccessToKen())
		if tk == nil {
			return nil, exception.NewUnauthorized("token required")
		}
		req.Namespace = tk.Namespace
	}

	descKey := pipeline.PipeLineObjectKey(req.Namespace, req.Id)
	i.log.Infof("describe etcd pipeline resource key: %s", descKey)
	resp, err := i.client.Get(ctx, descKey)
	if err != nil {
		return nil, err
	}

	if resp.Count == 0 {
		return nil, exception.NewNotFound("pipeline %s not found", req.Id)
	}

	if resp.Count > 1 {
		return nil, exception.NewInternalServerError("pipeline find more than one: %d", resp.Count)
	}

	ins := pipeline.NewDefaultPipeline()
	for index := range resp.Kvs {
		// 解析对象
		ins, err = pipeline.LoadPipelineFromBytes(resp.Kvs[index].Value)
		if err != nil {
			i.log.Error(err)
			continue
		}
		ins.ResourceVersion = resp.Header.Revision
	}
	return ins, nil
}

// DeletePipeline 删除时清除所有关联step
func (i *impl) DeletePipeline(ctx context.Context, req *pipeline.DeletePipelineRequest) (
	*pipeline.Pipeline, error) {
	in, err := gcontext.GetGrpcInCtx(ctx)
	if err != nil {
		return nil, err
	}
	tk := session.S().GetToken(in.GetAccessToKen())
	if tk == nil {
		return nil, exception.NewUnauthorized("token required")
	}

	ins, err := i.DescribePipeline(ctx, pipeline.NewDescribePipelineRequestWithID(req.Id))
	if err != nil {
		return nil, err
	}

	// 先删除pipeline对应的step
	if err := i.deletePipelineStep(ctx, ins); err != nil {
		i.log.Errorf("delete pipeline [%s] steps error, %s", err)
	}

	descKey := ins.MakeObjectKey()
	i.log.Infof("delete etcd pipeline resource key: %s", descKey)
	_, err = i.client.Delete(ctx, descKey, clientv3.WithPrevKV())
	if err != nil {
		return nil, err
	}

	return ins, nil
}

func (i *impl) WatchPipeline(*pipeline.WatchPipelineRequest, pipeline.Service_WatchPipelineServer) error {
	// 监听事件
	stepWatchKey := pipeline.EtcdStepPrefix()
	// i.watchChan = i.client.Watch(ctx, stepWatchKey, clientv3.WithPrefix())
	i.log.Infof("watch etcd step resource key: %s", stepWatchKey)
	return nil
}

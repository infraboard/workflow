package impl

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/infraboard/mcube/exception"
	"github.com/infraboard/mcube/grpc/gcontext"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc/metadata"

	"github.com/infraboard/workflow/api/apps/action"
	"github.com/infraboard/workflow/api/apps/pipeline"
)

func (i *impl) CreatePipeline(ctx context.Context, req *pipeline.CreatePipelineRequest) (
	*pipeline.Pipeline, error) {
	p, err := pipeline.NewPipeline(req)
	if err != nil {
		return nil, err
	}

	if err := i.validatePipelineStage(ctx, p); err != nil {
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

func (i *impl) validatePipelineStage(ctx context.Context, p *pipeline.Pipeline) error {
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
	_, err := i.action.DescribeAction(ctx, action.NewDescribeActionRequest(s.ActionName(), s.ActionVersion()))
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
	descReq := pipeline.NewDescribePipelineRequestWithID(req.Id)
	descReq.Namespace = req.Namespace
	ins, err := i.DescribePipeline(ctx, descReq)
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

func (i *impl) WatchPipeline(stream pipeline.Service_WatchPipelineServer) error {
	for {
		union, err := stream.Recv()
		if err != nil {
			return err
		}

		i.log.Debugf("get watch req: %s", union)

		req := union.GetCreateRequest()
		if req != nil {
			i.watchPipeline(req, stream)
			continue
		}

		cancel := union.GetCancelRequest()
		if cancel != nil {
			if fn, ok := i.watchCancel[cancel.WatchId]; ok {
				fn()
			}
			return nil
		}
	}
}

func (i *impl) watchPipeline(req *pipeline.CreateWatchPipelineRequest, stream pipeline.Service_WatchPipelineServer) error {
	if err := req.Validate(); err != nil {
		return exception.NewBadRequest("validate watch pipeline request error, %s", err)
	}

	opts := []clientv3.OpOption{}
	watchKey := pipeline.PipeLineObjectKey(req.Namespace, req.Id)
	switch req.Mod {
	case pipeline.PIPELINE_WATCH_MOD_BY_ID:
		opts = append(opts, clientv3.WithPrefix())
		// 检查pipeline是否存在
		if err := i.checkKeyIsExist(watchKey); err != nil {
			return exception.NewBadRequest("pipeline not found, %s", err)
		}
	case pipeline.PIPELINE_WATCH_MOD_BY_NAMESPACE:
		// 检查namespace是否存在
	default:
		return exception.NewBadRequest("unkwon watch mod %s", req.Mod)
	}

	if req.DryRun {
		i.log.Debug("dry run mode skip watch")
		return nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	watchChan := i.client.Watch(ctx, watchKey, opts...)
	watchId := i.SetWatcherCancelFn(cancel)

	i.log.Infof("watch etcd step resource key: %s, watch id %d start ...", watchKey, watchId)
	go i.dumpPipelineEvents(watchChan, stream, watchId)
	return nil
}

func (i *impl) dumpPipelineEvents(ch clientv3.WatchChan, stream pipeline.Service_WatchPipelineServer, watchId int64) {
	defer func() {
		i.log.Infof("watch : %d end ...", watchId)
	}()

	// 发送消息
	resp := &pipeline.WatchPipelineResponse{WatchId: watchId}
	if err := stream.Send(resp); err != nil {
		i.log.Errorf("send pipeline watch id error, %s", watchId, err)
		return
	}

	// 处理所有事件
	for ppResp := range ch {
		for index := range ppResp.Events {
			event := ppResp.Events[index]
			i.log.Debugf("receive pipeline event, %s", event.Kv.Key)
			// 解析对象
			ins, err := pipeline.LoadPipelineFromBytes(event.Kv.Value)
			if err != nil {
				i.sendFailed(stream, exception.NewInternalServerError("load pipeline from bytes error, %s", err))
				return
			}
			ins.ResourceVersion = ppResp.Header.Revision

			// 发送消息
			resp := &pipeline.WatchPipelineResponse{WatchId: watchId}
			resp.Pipeline = ins
			if err := stream.Send(resp); err != nil {
				i.log.Errorf("send pipeline %s events error, %s", ins.ShortDescribe(), err)
				return
			}
		}
	}
}

func (i *impl) sendFailed(stream pipeline.Service_WatchPipelineServer, err exception.APIException) {
	trailer := metadata.Pairs(
		gcontext.ResponseCodeHeader, fmt.Sprintf("%d", err.ErrorCode()),
		gcontext.ResponseReasonHeader, err.Reason(),
		gcontext.ResponseDescHeader, err.Error(),
	)
	stream.SetTrailer(trailer)
}

func (i *impl) checkKeyIsExist(key string) error {
	resp, err := i.client.Get(context.Background(), key)
	if err != nil {
		return err
	}

	if resp.Count == 0 {
		return exception.NewNotFound("key %s not found", key)
	}
	return nil
}

package impl

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/infraboard/keyauth/client/session"
	"github.com/infraboard/mcube/exception"
	"github.com/infraboard/mcube/grpc/gcontext"
	"github.com/rs/xid"
	clientv3 "go.etcd.io/etcd/client/v3"

	"github.com/infraboard/workflow/api/pkg/pipeline"
)

func (i *impl) CreateStep(ctx context.Context, req *pipeline.CreateStepRequest) (
	*pipeline.Step, error) {
	if err := req.Validate(); err != nil {
		return nil, exception.NewBadRequest("validate create step request error, %s", err)
	}

	step := pipeline.NewStep(pipeline.STEP_CREATE_BY_USER, req)
	step.Key = xid.New().String()

	if err := i.validateStep(ctx, step); err != nil {
		return nil, exception.NewBadRequest("validate step error, %s", err)
	}

	value, err := json.Marshal(step)
	if err != nil {
		return nil, err
	}

	objKey := step.MakeObjectKey()
	objValue := string(value)

	if _, err := i.client.Put(context.Background(), objKey, objValue); err != nil {
		return nil, fmt.Errorf("put step with key: %s, error, %s", objKey, err.Error())
	}
	i.log.Debugf("create step success, key: %s", objKey)
	return step, nil
}

func (i *impl) QueryStep(ctx context.Context, req *pipeline.QueryStepRequest) (
	*pipeline.StepSet, error) {
	listKey := pipeline.EtcdStepPrefix()
	i.log.Infof("list etcd step resource key: %s", listKey)
	resp, err := i.client.Get(ctx,
		listKey,
		clientv3.WithPrefix(),
		clientv3.WithSort(clientv3.SortByModRevision, clientv3.SortDescend))
	if err != nil {
		return nil, err
	}

	ps := pipeline.NewStepSet()
	for index := range resp.Kvs {
		// 解析对象
		ins, err := pipeline.LoadStepFromBytes(resp.Kvs[index].Value)
		if err != nil {
			i.log.Error(err)
			continue
		}
		ps.Add(ins)
	}
	return ps, nil
}

func (i *impl) DescribeStep(ctx context.Context, req *pipeline.DescribeStepRequest) (
	*pipeline.Step, error) {
	descKey := pipeline.StepObjectKey(req.Key)
	i.log.Infof("describe etcd step resource key: %s", descKey)
	resp, err := i.client.Get(ctx, descKey)
	if err != nil {
		return nil, err
	}

	if resp.Count == 0 {
		return nil, exception.NewNotFound("step %s not found", req.Key)
	}

	if resp.Count > 1 {
		return nil, exception.NewInternalServerError("step find more than one: %d", resp.Count)
	}

	ins := pipeline.NewDefaultStep()
	for index := range resp.Kvs {
		// 解析对象
		ins, err = pipeline.LoadStepFromBytes(resp.Kvs[index].Value)
		if err != nil {
			i.log.Error(err)
			continue
		}
	}
	return ins, nil
}

func (i *impl) DeleteStep(ctx context.Context, req *pipeline.DeleteStepRequest) (
	*pipeline.Step, error) {
	if err := req.Validate(); err != nil {
		return nil, exception.NewBadRequest("validate delete request error, %s", err)
	}
	descKey := pipeline.StepObjectKey(req.Key)
	i.log.Infof("delete etcd pipeline resource key: %s", descKey)
	resp, err := i.client.Delete(ctx, descKey, clientv3.WithPrevKV())
	if err != nil {
		return nil, err
	}

	if resp.Deleted == 0 {
		return nil, exception.NewNotFound("step %s not found", req.Key)
	}

	ins := pipeline.NewDefaultStep()
	for index := range resp.PrevKvs {
		// 解析对象
		ins, err = pipeline.LoadStepFromBytes(resp.PrevKvs[index].Value)
		if err != nil {
			i.log.Error(err)
			continue
		}
		ins.ResourceVersion = resp.Header.Revision
	}
	return ins, nil
}

func (i *impl) CancelStep(ctx context.Context, req *pipeline.CancelStepRequest) (
	*pipeline.Step, error) {
	in, err := gcontext.GetGrpcInCtx(ctx)
	if err != nil {
		return nil, err
	}

	message := ""
	tk := session.S().GetToken(in.GetAccessToKen())
	if tk == nil {
		message = fmt.Sprintf("account %s cancled step at %s", tk.Account, time.Now())
	}

	s, err := i.DescribeStep(ctx, pipeline.NewDescribeStepRequestWithKey(req.Key))
	if err != nil {
		return nil, err
	}

	s.Cancel(message)
	if err := i.putStep(ctx, s); err != nil {
		return nil, fmt.Errorf("update step error, %s", err)
	}

	return s, nil
}

func (i *impl) AuditStep(ctx context.Context, req *pipeline.AuditStepRequest) (
	*pipeline.Step, error) {
	s, err := i.DescribeStep(ctx, pipeline.NewDescribeStepRequestWithKey(req.Key))
	if err != nil {
		return nil, err
	}

	if !s.WithAudit {
		return nil, exception.NewBadRequest("this step needn't audit")
	}

	s.Audit(req.AuditReponse, req.AuditMessage)
	if err := i.putStep(ctx, s); err != nil {
		return nil, fmt.Errorf("update step error, %s", err)
	}

	return s, nil
}

func (i *impl) putStep(ctx context.Context, ins *pipeline.Step) error {
	value, err := json.Marshal(ins)
	if err != nil {
		return err
	}

	objKey := ins.MakeObjectKey()
	objValue := string(value)

	if _, err := i.client.Put(context.Background(), objKey, objValue); err != nil {
		return fmt.Errorf("put step with key: %s, error, %s", objKey, err.Error())
	}
	i.log.Debugf("put step success, key: %s", objKey)

	return nil
}

func (i *impl) deletePipelineStep(ctx context.Context, ins *pipeline.Pipeline) error {
	prefix := ins.StepPrefix()
	if len(prefix) < 12 {
		return fmt.Errorf("prefix length must large than 12")
	}

	deletePrefixKey := pipeline.StepObjectKey(prefix)
	i.log.Infof("delete etcd step resource key prefix: %s", deletePrefixKey)
	resp, err := i.client.Delete(ctx, deletePrefixKey, clientv3.WithPrefix())
	if err != nil {
		return err
	}

	i.log.Infof("delete pipeline %s total %d steps", ins.ShortDescribe(), resp.Deleted)
	return nil
}

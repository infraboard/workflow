package etcd

import (
	"context"
	"errors"
	"strings"

	"github.com/infraboard/mcube/logger"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"

	"github.com/infraboard/workflow/api/pkg/node"
	"github.com/infraboard/workflow/api/pkg/pipeline"
	"github.com/infraboard/workflow/scheduler/informer"
)

type shared struct {
	log               logger.Logger
	client            clientv3.Watcher
	nodeHandler       informer.NodeEventHandler
	pipelineHandler   informer.PipelineEventHandler
	stepHandler       informer.StepEventHandler
	nodeWatchChan     clientv3.WatchChan
	pipelineWatchChan clientv3.WatchChan
	stepWatchChan     clientv3.WatchChan
}

// AddEventHandler 添加事件处理回调
func (i *shared) AddNodeEventHandler(h informer.NodeEventHandler) {
	i.nodeHandler = h
}

// AddPipelineEventHandler 添加事件处理回调
func (i *shared) AddPipelineTaskEventHandler(h informer.PipelineEventHandler) {
	i.pipelineHandler = h
}

func (i *shared) AddStepEventHandler(h informer.StepEventHandler) {
	i.stepHandler = h
}

// Run 启动 Watch
func (i *shared) Run(ctx context.Context) error {
	// 是否准备完成
	if err := i.isReady(); err != nil {
		return err
	}

	// 监听事件
	i.watchAll(ctx)

	// 后台处理事件
	go i.dealEvents()
	return nil
}

func (i *shared) dealEvents() {
	// 处理所有事件
	for {
		select {
		case nodeResp := <-i.nodeWatchChan:
			for _, event := range nodeResp.Events {
				if err := i.notifyNode(event, nodeResp.Header.GetRevision()); err != nil {
					i.log.Error(err)
				}
			}
		case ppResp := <-i.pipelineWatchChan:
			for _, event := range ppResp.Events {
				if err := i.notifyPipeline(event, ppResp.Header.GetRevision()); err != nil {
					i.log.Error(err)
				}
			}
		case stepResp := <-i.stepWatchChan:
			for _, event := range stepResp.Events {
				if err := i.notifyStep(event, stepResp.Header.GetRevision()); err != nil {
					i.log.Error(err)
				}
			}
		}
	}
}

func (i *shared) isReady() error {
	if i.nodeHandler == nil {
		return errors.New("NodeEventHandler not add")
	}
	if i.pipelineHandler == nil {
		return errors.New("PipelineEventHandler not add")
	}
	if i.stepHandler == nil {
		return errors.New("StepEventHandler not add")
	}
	return nil
}

func (i *shared) watchAll(ctx context.Context) {
	// 监听事件
	nodeWatchKey := node.EtcdNodePrefixWithType(node.NodeType)
	i.nodeWatchChan = i.client.Watch(ctx, nodeWatchKey, clientv3.WithPrefix())
	i.log.Infof("watch etcd node resource key: %s", nodeWatchKey)

	ppWatchKey := pipeline.EtcdPipelinePrefix()
	i.pipelineWatchChan = i.client.Watch(ctx, ppWatchKey, clientv3.WithPrefix())
	i.log.Infof("watch etcd pipeline resource key: %s", ppWatchKey)

	stepWatchKey := pipeline.EtcdStepPrefix()
	i.stepWatchChan = i.client.Watch(ctx, stepWatchKey, clientv3.WithPrefix())
	i.log.Infof("watch etcd step resource key: %s", stepWatchKey)
}

func (i *shared) notifyNode(event *clientv3.Event, eventVersion int64) error {
	// 解析对象
	obj, err := node.LoadNodeFromBytes(event.Kv.Value)
	if err != nil {
		return err
	}

	// 解析事件为删除事件的时候，value内容为空，从key中解析 node 的 serviceName and instanceName
	if obj == nil {
		obj = new(node.Node)
		etcdNodeKeyList := strings.Split(string(event.Kv.Key), "/")
		obj.InstanceName = etcdNodeKeyList[len(etcdNodeKeyList)-1]
		obj.ServiceName = etcdNodeKeyList[1]
	}
	obj.ResourceVersion = eventVersion
	switch event.Type {
	case mvccpb.PUT:
		i.nodeHandler.OnAdd(obj)
	case mvccpb.DELETE:
		i.nodeHandler.OnDelete(obj)
	default:
	}
	return nil
}

func (i *shared) notifyPipeline(event *clientv3.Event, eventVersion int64) error {
	// 解析对象
	obj, err := pipeline.LoadPipelineFromBytes(event.Kv.Value)
	if err != nil {
		return err
	}
	switch event.Type {
	case mvccpb.PUT:
		i.pipelineHandler.OnAdd(obj)
	case mvccpb.DELETE:
		i.pipelineHandler.OnDelete(obj)
	default:
	}
	return nil
}

func (i *shared) notifyStep(event *clientv3.Event, eventVersion int64) error {
	// 解析对象
	obj, err := pipeline.LoadStepFromBytes(event.Kv.Value)
	if err != nil {
		return err
	}
	switch event.Type {
	case mvccpb.PUT:
		i.stepHandler.OnAdd(obj)
	case mvccpb.DELETE:
		i.stepHandler.OnDelete(obj)
	default:
	}
	return nil
}

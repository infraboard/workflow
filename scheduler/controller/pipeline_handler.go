package controller

import (
	"context"
	"time"

	"github.com/infraboard/workflow/api/pkg/node"
	"github.com/infraboard/workflow/api/pkg/pipeline"
	"github.com/infraboard/workflow/scheduler/informer"
)

// 添加节点后, 需要延迟扫描, 从新调度未调度的Pipeline
func (c *PipelineScheduler) addNode(n *node.Node) {
	c.nodes.AddNode(n)
}
func (c *PipelineScheduler) delNode(n *node.Node) {
	c.nodes.DelNode(n)

	// 1. 获取该节点上所有pipeline
	ctx, cancle := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancle()
	ps, err := c.lister.ListPipeline(ctx, &informer.QueryPipelineOptions{Node: n.InstanceName})
	if err != nil {
		c.log.Errorf("list node pipelines error, %s", err)
	}

	// 2. 发送给调度队列
	for i := range ps.Items {
		c.workqueue.Add(ps.Items[i])
	}
}

func (c *PipelineScheduler) addPipeline(p *pipeline.Pipeline) {
	c.log.Infof("receive add object: %s", p)
	if err := p.Validate(); err != nil {
		c.log.Errorf("invalidate node error, %s", err)
		return
	}

	// 未调度的交给调度
	if len(p.ScheduledNodeName()) == 0 {
		c.workqueue.AddRateLimited(p)
	}
}

func (c *PipelineScheduler) addStep(p *pipeline.Step) {
	c.log.Infof("receive add object: %s", p)
	if err := p.Validate(); err != nil {
		c.log.Errorf("invalidate node error, %s", err)
		return
	}

	// 未调度的交给调度
}

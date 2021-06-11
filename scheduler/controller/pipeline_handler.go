package controller

import (
	"context"
	"time"

	"github.com/infraboard/workflow/api/pkg/node"
	"github.com/infraboard/workflow/api/pkg/pipeline"
	"github.com/infraboard/workflow/api/pkg/task"
	"github.com/infraboard/workflow/scheduler/informer"
)

// 添加节点后, 需要延迟扫描, 从新调度未调度的Pipeline
func (c *PipelineTaskScheduler) addNode(n *node.Node) {
	c.nodes.AddNode(n)
}
func (c *PipelineTaskScheduler) delNode(n *node.Node) {
	c.nodes.DelNode(n)

	// 1. 获取该节点上所有pipeline
	ctx, cancle := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancle()
	ps, err := c.lister.ListPipelineTask(ctx, &informer.QueryPipelineTaskOptions{Node: n.InstanceName})
	if err != nil {
		c.log.Errorf("list node pipelines error, %s", err)
	}

	// 2. 发送给调度队列
	for i := range ps.Items {
		c.workqueue.Add(ps.Items[i])
	}
}

// 每添加一个pipeline task
func (c *PipelineTaskScheduler) addPipelineTask(t *task.PipelineTask) {
	c.log.Infof("receive add object: %s", t)
	if err := t.Validate(); err != nil {
		c.log.Errorf("invalidate node error, %s", err)
		return
	}

	// 已经被处理的task无效再处理
	if t.SchedulerNodeName() != "" {
		c.log.Infof("pipeline task %s has scheduler node: %s", t.SchedulerNodeName())
		return
	}

	// TODO: 使用分布式锁trylock处 理多个实例竞争调度问题

	// 将需要调度的任务, 交给step调度器调度
	steps := t.NextStep()
	for i := range steps {
		step := steps[i]
		c.log.Debugf("add pipeline step: %s", step.Key)
		err := c.lister.UpdateStep(step)
		if err != nil {
			c.log.Errorf(err.Error())
		}
	}
}

func (c *PipelineTaskScheduler) addStep(s *pipeline.Step) {
	c.log.Infof("receive add object: %s", s)
	if err := s.Validate(); err != nil {
		c.log.Errorf("invalidate node error, %s", err)
		return
	}

	// 已经调度的任务不处理
	nn := s.ScheduledNodeName()
	if nn != "" {
		c.log.Infof("step %s has scheuled to node %s", s.Key, nn)
	}

	c.workqueue.AddRateLimited(s)
}

func (c *PipelineTaskScheduler) deleteStep(p *pipeline.Step) {
	c.log.Infof("receive add object: %s", p)
	if err := p.Validate(); err != nil {
		c.log.Errorf("invalidate node error, %s", err)
		return
	}

	// 未调度的交给调度
}

func (c *PipelineTaskScheduler) updateStep(oldObj, newObj *pipeline.Step) {

}

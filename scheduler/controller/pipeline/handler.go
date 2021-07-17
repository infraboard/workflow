package pipeline

import (
	"context"
	"time"

	"github.com/infraboard/workflow/api/pkg/node"
	"github.com/infraboard/workflow/api/pkg/pipeline"
)

// 添加节点后, 需要延迟扫描, 从新调度未调度的Pipeline
func (c *PipelineScheduler) addNode(n *node.Node) {
	c.nodes.Add(n)
}
func (c *PipelineScheduler) delNode(n *node.Node) {
	c.nodes.Delete(n)

	// 1. 获取该节点上所有pipeline
	ctx, cancle := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancle()
	ps, err := c.pi.Lister().List(ctx, &pipeline.QueryPipelineOptions{Node: n.InstanceName})
	if err != nil {
		c.log.Errorf("list node pipelines error, %s", err)
	}

	// 2. 发送给调度队列
	for i := range ps.Items {
		c.workqueue.Add(ps.Items[i])
	}
}

// 每添加一个pipeline
func (c *PipelineScheduler) addPipeline(p *pipeline.Pipeline) {
	c.log.Debugf("[pipeline] receive add pipeline: %s status: %s", p.ShortDescribe(), p.Status.Status)
	if err := p.Validate(); err != nil {
		c.log.Errorf("invalidate pipeline error, %s", err)
		return
	}

	// 已经处理完成的无需处理
	if p.IsComplete() {
		c.log.Debugf("skip run complete pipeline %s, status: %s", p.ShortDescribe(), p.Status.Status)
		return
	}

	// TODO: 使用分布式锁trylock处理 多个实例竞争调度问题

	// 标记开始执行, 并更新保存
	if !p.Status.IsRunning() {
		p.Status.Run()
		if err := c.pi.Recorder().Update(p); err != nil {
			c.log.Errorf("update pipeline %s status to store error, %s", p.ShortDescribe(), err)
		}
		return
	}

	// 将需要调度的任务, 交给step调度器调度
	steps := p.NextStep()
	c.log.Debugf("pipeline %s next step is %v", p.ShortDescribe(), steps)
	for i := range steps {
		step := steps[i]
		c.log.Debugf("create pipeline step: %s", step.Key)
		err := c.si.Recorder().Update(step)
		if err != nil {
			c.log.Errorf(err.Error())
		}
	}
}

// func (c *PipelineScheduler) addStep(s *pipeline.Step) {
// 	c.log.Infof("[step] receive add step: %s", s)
// 	if err := s.Validate(); err != nil {
// 		c.log.Errorf("invalidate node error, %s", err)
// 		return
// 	}

// 	// 已经调度的任务不处理
// 	nn := s.ScheduledNodeName()
// 	if nn != "" {
// 		c.log.Infof("step %s has scheuled to node %s", s.Key, nn)
// 	}

// 	// 判断是否需要审批, 审批通过后放行
// 	if s.WithAudit && !s.IsAudit() {
// 		// TODO: 发送审批事件
// 		s.Status.Status = pipeline.STEP_STATUS_AUDITING
// 		c.log.Debug("send audit notify")
// 	}

// 	c.workqueue.AddRateLimited(s)
// }

// func (c *PipelineScheduler) deleteStep(p *pipeline.Step) {
// 	c.log.Infof("receive add object: %s", p)
// 	if err := p.Validate(); err != nil {
// 		c.log.Errorf("invalidate node error, %s", err)
// 		return
// 	}

// 	// 未调度的交给调度
// }

// // 如果step有状态更新, 判断step是否执行结束
// // 如果结束就将step的状态同步更新到Pipeline上去, 再删除step
// func (c *PipelineScheduler) updateStep(oldObj, newObj *pipeline.Step) {
// 	if !newObj.IsComplete() {
// 		c.log.Debugf("step %s status is %s, skip sync to pipeline", newObj.Status.Status)
// 		return
// 	}

// 	c.log.Debugf("step %s status [%s] is complete, start sync to pipeline ...")

// }

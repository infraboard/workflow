package pipeline

import (
	"context"
	"errors"
	"fmt"

	"github.com/infraboard/workflow/api/pkg/pipeline"
)

// syncHandler compares the actual state with the desired, and attempts to
// converge the two. It then updates the Status block of the Network resource
// with the current status of the resource.
func (c *Controller) syncHandler(key string) error {
	obj, ok, err := c.informer.GetStore().GetByKey(key)
	if err != nil {
		return err
	}

	ins, isOK := obj.(*pipeline.Pipeline)
	if !isOK {
		return errors.New("invalidate *pipeline.Pipeline obj")
	}

	// 如果不存在, 这期望行为为删除 (DEL)
	if !ok {
		c.log.Debugf("wating remove step: %s", key)
		if err := c.deletePipeline(ins); err != nil {
			return err
		}
		c.log.Infof("remove success, step: %s", key)
	}

	// 运行pipeline
	if err := c.runPipeline(ins); err != nil {
		return err
	}

	return nil
}

func (c *Controller) deletePipeline(*pipeline.Pipeline) error {
	c.log.Debugf("sync delete pipeline ...")
	return nil
}

// 运行一个pipeline，流程如下
// 1. 判断是否已经完成, 已完不做处理
// 2. 判断是否需要调度, 未调度先调度
// 3. 判断是否已经运行, 未运行先标记为运行状态
// 4. 如果是运行状态, 判断pipline是否需要中断
// 5. 如果pipeline正常, 则开始运行定义的 Next Step
func (c *Controller) runPipeline(p *pipeline.Pipeline) error {
	c.log.Debugf("receive add pipeline: %s status: %s", p.ShortDescribe(), p.Status.Status)
	if err := p.Validate(); err != nil {
		return fmt.Errorf("invalidate pipeline error, %s", err)
	}

	// 已经处理完成的无需处理
	if p.IsComplete() {
		return fmt.Errorf("skip run complete pipeline %s, status: %s", p.ShortDescribe(), p.Status.Status)
	}

	// TODO: 使用分布式锁trylock处理 多个实例竞争调度问题

	// 未调度的选进行调度后, 再处理
	if !p.IsScheduled() {
		if err := c.schedulePipeline(p); err != nil {
			return err
		}
		return nil
	}

	// 标记开始执行, 并更新保存
	if !p.IsRunning() {
		p.Run()
		if err := c.informer.Recorder().Update(p); err != nil {
			c.log.Errorf("update pipeline %s start status to store error, %s", p.ShortDescribe(), err)
		} else {
			c.log.Debugf("update pipeline %s start status to store success", p.ShortDescribe())
		}
		return nil
	}

	// 判断pipeline没有要执行的下一步, 则结束整个Pipeline
	steps := c.nextStep(p)
	c.log.Debugf("pipeline %s start run next steps: %s", p.ShortDescribe(), steps)
	return c.runPipelineNextStep(steps)
}

func (c *Controller) nextStep(p *pipeline.Pipeline) []*pipeline.Step {
	// 取消 pipeline 下次执行需要的step
	steps, ok := p.NextStep()
	if ok {
		if err := c.informer.Recorder().Update(p); err != nil {
			c.log.Errorf("update pipeline %s end status to store error, %s", p.ShortDescribe(), err)
		} else {
			c.log.Debugf("pipeline is complete, update pipeline status to db success")
		}
		return nil
	}

	// 找出需要同步的step
	needSync := []*pipeline.Step{}
	for i := range steps {
		ins := steps[i]

		// 判断step是否已经运行, 如果已经运行则更新Pipeline状态
		old, err := c.step.Lister().Get(context.Background(), ins.Key)
		if err != nil {
			c.log.Errorf("get step %s by key error, %s", ins.Key, err)
			return nil
		}

		if old == nil {
			c.log.Debugf("step %s not found in db", ins.Key)
			continue
		}

		// 状态相等 则无需同步
		if ins.Status.Status.Equal(old.Status.Status) {
			c.log.Debugf("pipeline step status: %s, etcd step status: %s, has sync",
				ins.Status.Status, old.Status.Status)
			continue
		}

		needSync = append(needSync, old)
	}

	// 同步step到pipeline上
	if len(needSync) > 0 {
		for i := range needSync {
			c.log.Debugf("sync step %s to pipeline ...", needSync[i].Key)
			p.UpdateStep(needSync[i])
		}
		if err := c.informer.Recorder().Update(p); err != nil {
			c.log.Errorf("update pipeline status error, %s", err)
			return nil
		}
		c.log.Debugf("sync %d steps ok", len(needSync))
		return nil
	}

	return steps
}

func (c *Controller) runPipelineNextStep(steps []*pipeline.Step) error {
	// 将需要调度的任务, 交给step调度器调度
	if c.step == nil {
		return fmt.Errorf("step recorder is nil")
	}

	// 有step则进行执行
	for i := range steps {
		ins := steps[i]

		// 如果开启审核，需要通过后，才能执行
		if !c.isAllow(ins) {
			continue
		}

		c.log.Debugf("create pipeline step: %s", ins.Key)
		if err := c.step.Recorder().Update(ins); err != nil {
			c.log.Errorf(err.Error())
		}
	}

	return nil
}

func (c *Controller) isAllow(s *pipeline.Step) bool {
	if !s.WithAudit {
		return true
	}

	// TODO: 如果未处理, 发送通知
	if !s.HasSendAuditNotify() {
		// TODO:
		c.log.Errorf("send notify ...")
		s.MarkSendAuditNotify()
		// 更新step
	}

	if s.AuditPass() {
		c.log.Debugf("step %s waiting for audit", s.Key)
		return true
	}

	return false
}

// Pipeline 调度
func (c *Controller) schedulePipeline(p *pipeline.Pipeline) error {
	c.log.Debugf("pipeline %s start schedule ...", p.ShortDescribe())

	node, err := c.picker.Pick(p)
	if err != nil {
		return err
	}

	// 没有合法的node
	if node == nil {
		return fmt.Errorf("no excutable scheduler")
	}

	c.log.Debugf("choice scheduler %s for pipeline %s", node.InstanceName, p.Id)
	p.SetScheduleNode(node.InstanceName)
	c.updatePipelineStatus(p)
	return nil
}

func (c *Controller) updatePipelineStatus(p *pipeline.Pipeline) {
	if p == nil {
		c.log.Errorf("update pipeline is nil")
		return
	}

	if c.informer.Recorder() == nil {
		c.log.Errorf("pipeline informer recorder missed")
		return
	}

	// 清除一下其他数据
	if err := c.informer.Recorder().Update(p); err != nil {
		c.log.Errorf("update scheduled pipeline error, %s", err)
	}
}

// step 如果完成后, 将状态记录到Pipeline上, 并删除step
func (c *Controller) UpdateStepCallback(old, new *pipeline.Step) {
	c.log.Debugf("receive step update event, start update step status to pipeline ...")

	if !new.IsComplete() {
		c.log.Debugf("step status is %s, skip update to pipeline", new.Status.Status)
		return
	}

	key := pipeline.PipeLineObjectKey(new.GetNamespace(), new.GetPipelineID())
	obj, ok, err := c.informer.GetStore().GetByKey(key)
	if err != nil {
		c.log.Errorf("get pipeline from store error, %s", err)
		return
	}

	if !ok {
		c.log.Errorf("pipeline key %s not found in store cache", key)
		return
	}

	p, isOK := obj.(*pipeline.Pipeline)
	if !isOK {
		c.log.Errorf("invalidate *pipeline.Pipeline obj")
		return
	}

	if err := p.UpdateStep(new); err != nil {
		c.log.Errorf("update pipeline step error, %s", err)
		return
	}

	if err := c.informer.Recorder().Update(p); err != nil {
		c.log.Errorf("update pipeline status to store error, %s", err)
		return
	}

	c.log.Debugf("update pipeline %s step %s success", p.ShortDescribe(), new.Key)
}

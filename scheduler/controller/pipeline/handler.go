package pipeline

import (
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

	st, isOK := obj.(*pipeline.Pipeline)
	if !isOK {
		return errors.New("invalidate *pipeline.Pipeline obj")
	}

	// 如果不存在, 这期望行为为删除 (DEL)
	if !ok {
		c.log.Debugf("wating remove step: %s", key)
		if err := c.deletePipeline(st); err != nil {
			return err
		}
		c.log.Infof("remove success, step: %s", key)
	}

	// 添加
	if err := c.addPipeline(st); err != nil {
		return err
	}

	return nil
}

func (c *Controller) deletePipeline(*pipeline.Pipeline) error {
	c.log.Debugf("sync delete pipeline ...")
	return nil
}

// 每添加一个pipeline
func (c *Controller) addPipeline(p *pipeline.Pipeline) error {
	c.log.Debugf("receive add pipeline: %s status: %s", p.ShortDescribe(), p.Status.Status)
	if err := p.Validate(); err != nil {
		return fmt.Errorf("invalidate pipeline error, %s", err)
	}

	// 已经处理完成的无需处理
	if p.IsComplete() {
		return fmt.Errorf("skip run complete pipeline %s, status: %s", p.ShortDescribe(), p.Status.Status)
	}

	// TODO: 使用分布式锁trylock处理 多个实例竞争调度问题

	// 执行调度
	if err := c.schedulePipeline(p); err != nil {
		return err
	}

	// 标记开始执行, 并更新保存
	if !p.IsRunning() {
		p.Run()
		if err := c.informer.Recorder().Update(p); err != nil {
			c.log.Errorf("update pipeline %s status to store error, %s", p.ShortDescribe(), err)
		}
		return nil
	}

	// 将需要调度的任务, 交给step调度器调度
	if c.step == nil {
		return fmt.Errorf("step recorder is nil")
	}

	steps := p.NextStep()
	c.log.Debugf("pipeline %s next step is %v", p.ShortDescribe(), steps)
	for i := range steps {
		step := steps[i]
		c.log.Debugf("create pipeline step: %s", step.Key)
		if err := c.step.Recorder().Update(step); err != nil {
			c.log.Errorf(err.Error())
		}
	}

	return nil
}

// Pipeline 调度
func (c *Controller) schedulePipeline(p *pipeline.Pipeline) error {
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
func (c *Controller) StepUpdate(old, new *pipeline.Step) {
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

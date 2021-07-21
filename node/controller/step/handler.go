package step

import (
	"errors"

	"github.com/infraboard/workflow/api/pkg/pipeline"
	"github.com/infraboard/workflow/node/controller/step/engine"
)

// syncHandler compares the actual state with the desired, and attempts to
// converge the two. It then updates the Status block of the Network resource
// with the current status of the resource.
func (c *Controller) syncHandler(key string) error {
	obj, ok, err := c.informer.GetStore().GetByKey(key)
	if err != nil {
		return err
	}

	st, isOK := obj.(*pipeline.Step)
	if !isOK {
		return errors.New("invalidate *pipeline.Step obj")
	}

	// 如果不存在, 这期望行为为删除 (DEL)
	if !ok {
		c.log.Debugf("wating remove step: %s", key)
		if err := c.deleteStep(st); err != nil {
			return err
		}
		c.log.Infof("remove success, step: %s", key)
	}

	// 添加step
	if err := c.addStep(st); err != nil {
		return err
	}
	return nil
}

func (c *Controller) addStep(s *pipeline.Step) error {
	// 开始执行, 更新状态
	s.Run()
	c.updateStepStatus(s)

	// 执行结束, 更新状态
	engine.RunStep(s)
	c.updateStepStatus(s)
	return nil
}

func (c *Controller) updateStepStatus(s *pipeline.Step) {
	if err := c.informer.Recorder().Update(s); err != nil {
		c.log.Errorf("update step end %s status error, %s", s.Key, err)
	}
	c.log.Debugf("update step status: %s", s.Status)
}

func (c *Controller) deleteStep(p *pipeline.Step) error {
	c.log.Infof("receive add object: %s", p)
	if err := p.Validate(); err != nil {
		c.log.Errorf("invalidate node error, %s", err)
		return nil
	}

	//

	// 未调度的交给调度
	return nil
}

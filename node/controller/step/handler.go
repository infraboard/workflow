package step

import (
	"errors"
	"fmt"

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
	status := s.Status.Status
	switch status {
	case pipeline.STEP_STATUS_PENDDING:
		engine.RunStep(s)
		return nil
	case pipeline.STEP_STATUS_RUNNING:
		// TODO: 判断引擎中该step状态是否一致
		// 如果不一致则同步状态, 但是不作再次运行
		c.log.Debugf("step is running, no thing todo")
	case pipeline.STEP_STATUS_CANCELING:
		return c.cancelStep(s)
	case pipeline.STEP_STATUS_SUCCEEDED,
		pipeline.STEP_STATUS_FAILED,
		pipeline.STEP_STATUS_CANCELED,
		pipeline.STEP_STATUS_SKIP,
		pipeline.STEP_STATUS_REFUSE:
		return fmt.Errorf("step %s status is %s has complete", s.Key, status)
	case pipeline.STEP_STATUS_AUDITING:
		return fmt.Errorf("step %s is %s, is auditing", s.Key, status)
	}

	return nil
}

func (c *Controller) cancelStep(s *pipeline.Step) error {
	return nil
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

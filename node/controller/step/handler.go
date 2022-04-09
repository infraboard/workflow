package step

import (
	"context"
	"errors"
	"fmt"

	"github.com/infraboard/workflow/api/apps/pipeline"
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

	// 如果不存在, 这期望行为为删除 (DEL)
	if !ok {
		c.log.Debugf("remove step: %s, skip", key)
	}

	st, isOK := obj.(*pipeline.Step)
	if !isOK {
		return errors.New("invalidate *pipeline.Step obj")
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
		engine.RunStep(context.Background(), s)
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
	c.log.Infof("receive cancel object: %s", s)
	if err := s.Validate(); err != nil {
		c.log.Errorf("invalidate node error, %s", err)
		return nil
	}

	// 已经完成的step不作处理
	if s.IsComplete() {
		c.log.Debugf("step [%s] is complete, skip cancel", s.Key)
	}

	engine.CancelStep(s)
	return nil
}

// 当step删除时, 如果任务还在运行, 直接kill掉该任务
func (c *Controller) deleteStep(key string) error {
	// 取消任务

	return nil
}

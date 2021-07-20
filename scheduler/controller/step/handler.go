package step

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

	// 添加
	if err := c.addStep(st); err != nil {
		return err
	}

	return nil
}

func (c *Controller) addStep(s *pipeline.Step) error {
	c.log.Infof("receive add step: %s", s)
	if err := s.Validate(); err != nil {
		return fmt.Errorf("invalidate node error, %s", err)
	}

	// 已经调度的任务不处理
	nn := s.ScheduledNodeName()
	if nn != "" {
		c.log.Infof("step %s has scheuled to node %s", s.Key, nn)
	}

	// 判断是否需要审批, 审批通过后放行
	if s.WithAudit && !s.IsAudit() {
		// TODO: 发送审批事件
		s.Status.Status = pipeline.STEP_STATUS_AUDITING
		c.log.Debug("send audit notify")
	}

	if err := c.scheduleStep(s); err != nil {
		return err
	}

	return nil
}

// Step任务调度
func (c *Controller) scheduleStep(step *pipeline.Step) error {
	node, err := c.picker.Pick(step)
	if err != nil {
		return err
	}

	// 没有合法的node
	if node == nil {
		return fmt.Errorf("no available nodes")
	}

	c.log.Debugf("choice [%s] %s for step %s", node.Type, node.InstanceName, step.Key)
	step.SetScheduleNode(node.InstanceName)
	// 清除一下其他数据
	if err := c.informer.Recorder().Update(step); err != nil {
		c.log.Errorf("update scheduled step error, %s", err)
	}
	return nil
}

func (c *Controller) deleteStep(p *pipeline.Step) error {
	c.log.Infof("receive add object: %s", p)
	if err := p.Validate(); err != nil {
		c.log.Errorf("invalidate node error, %s", err)
		return nil
	}

	// 未调度的交给调度
	return nil
}

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
	if s.IsScheduled() {
		return fmt.Errorf("step %s has scheuled to node %s", s.Key, s.ScheduledNodeName())
	}

	// 如果开启审核，需要通过后，才能调度执行
	if !c.isAllow(s) {
		return fmt.Errorf("step not allow")
	}

	if err := c.scheduleStep(s); err != nil {
		return err
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

	// 审核通过 允许执行
	if s.AuditPass() {
		c.log.Debugf("step %s waiting for audit", s.Key)
		return true
	}

	return false
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
	if err := c.informer.Recorder().Update(step.Clone()); err != nil {
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

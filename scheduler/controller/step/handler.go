package step

import (
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

	// 如果不存在, 这期望行为为删除 (DEL)
	if !ok {
		c.log.Debugf("removed step: %s, skip", key)
		return nil
	}

	st, isOK := obj.(*pipeline.Step)
	if !isOK {
		return fmt.Errorf("object %T invalidate, is not *pipeline.Step obj, ", obj)
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
		return fmt.Errorf("step %s has schedule to node %s, skip add", s.Key, s.ScheduledNodeName())
	}

	// 已经调度的任务不处理
	if s.IsScheduledFailed() {
		return fmt.Errorf("step %s schedule failed, skip add", s.Key)
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

	// 判断审批状态是否同步

	// TODO: 如果未处理, 发送通知
	if !s.HasSendAuditNotify() {
		// TODO:
		c.log.Errorf("send notify ...")
		s.MarkSendAuditNotify()
		// 更新step
		if err := c.informer.Recorder().Update(s.Clone()); err != nil {
			c.log.Errorf("update scheduled step to auditing error, %s", err)
		}
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
	if err != nil || node == nil {
		c.log.Warnf("step %s pick node error, %s", step.Name, err)
		step.ScheduleFailed(err.Error())
		// 清除一下其他数据
		if err := c.informer.Recorder().Update(step.Clone()); err != nil {
			c.log.Errorf("update scheduled step error, %s", err)
		}
		return err
	}

	c.log.Debugf("choice [%s] %s for step %s", node.Type, node.InstanceName, step.Key)
	step.SetScheduleNode(node.InstanceName)
	// 清除一下其他数据
	if err := c.informer.Recorder().Update(step.Clone()); err != nil {
		c.log.Errorf("update scheduled step error, %s", err)
	}
	return nil
}

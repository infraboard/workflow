package step

import (
	"github.com/infraboard/workflow/api/pkg/pipeline"
)

func (c *Controller) addStep(s *pipeline.Step) {
	c.log.Infof("[step] receive add step: %s", s)
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

func (c *Controller) deleteStep(p *pipeline.Step) {
	c.log.Infof("receive add object: %s", p)
	if err := p.Validate(); err != nil {
		c.log.Errorf("invalidate node error, %s", err)
		return
	}

	// 未调度的交给调度
}

// 如果step有状态更新, 同步更新到Pipeline上去
func (c *Controller) updateStep(oldObj, newObj *pipeline.Step) {

}

package step

import (
	"errors"

	"github.com/infraboard/workflow/api/pkg/pipeline"
	"github.com/infraboard/workflow/node/controller/step/runner"
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
		if err := c.expectDelete(st); err != nil {
			return err
		}
		c.log.Infof("remove success, step: %s", key)
	}

	// 如果存在, 这期望行为为更新 (Update for DEL)
	// if c.cronPool.IsJobExist(job.HashID()) {
	// 	if err := c.cronPool.RemoveJob(job.HashID()); err != nil {
	// 		c.log.Error(err)
	// 	} else {
	// 		c.log.Infof("成功移除Cron(%s): %s.%s", strings.TrimSpace(job.HashID()), job.ProviderName, job.ExcutorName)
	// 	}
	// }

	return runner.RunStep(st)
}

func (c *Controller) expectDelete(s *pipeline.Step) error {
	// j, err := informer.NewJobFromStoreKey(key)
	// if err != nil {
	// 	return err
	// }
	// c.cronPool.RemoveJob(j.HashID())

	return nil
}

func (c *Controller) addStep(s *pipeline.Step) {
	c.log.Infof("receive add step: %s", s)
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

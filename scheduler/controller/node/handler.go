package node

import (
	"context"
	"fmt"

	"github.com/infraboard/workflow/api/apps/node"
)

// syncHandler compares the actual state with the desired, and attempts to
// converge the two. It then updates the Status block of the Network resource
// with the current status of the resource.
func (c *Controller) syncHandler(key string) error {
	c.log.Debugf("sync key: %s", key)

	obj, ok, err := c.informer.GetStore().GetByKey(key)
	if err != nil {
		return err
	}

	// 如果不存在, 这期望行为为删除 (DEL)
	if !ok {
		c.log.Debugf("remove node: %s, skip", key)
		return nil
	}

	n, isOK := obj.(*node.Node)
	if !isOK {
		return fmt.Errorf("object %T invalidate, is not *node.Node obj, ", obj)
	}

	return c.HandleAdd(n)
}

// 当有新的节点加入时, 那些调度失败的节点需要重新调度
func (c *Controller) HandleAdd(n *node.Node) error {
	// 补充重新调度的逻辑
	steps, err := c.stepLister.List(context.Background())
	if err != nil {
		c.log.Errorf("list steps error, %s", err)
		return nil
	}

	// 该删除节点上运行中的step进行重新调度
	for i := range steps {
		s := steps[i]
		if s.IsScheduledFailed() {
			c.log.Infof("step %s schedule failed, need reschedule ...", s.Key)
			s.SetScheduleNode("")
			err := c.stepRecorder.Update(s)
			if err != nil {
				c.log.Errorf("update step for reschedule error, %s", err)
				continue
			}
			c.log.Infof("reset step %s schedule node to \"\", waiting for reschedule", s.Key)
		}
	}
	return nil
}

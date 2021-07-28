package step

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/logger/zap"
	"k8s.io/client-go/util/workqueue"

	"github.com/infraboard/workflow/api/pkg/pipeline"
	"github.com/infraboard/workflow/common/cache"
	"github.com/infraboard/workflow/common/hooks"
	"github.com/infraboard/workflow/common/informers/step"
	"github.com/infraboard/workflow/scheduler/algorithm"
	"github.com/infraboard/workflow/scheduler/algorithm/roundrobin"
)

// NewStepController pipeline controller
func NewStepController(
	schedulerName string,
	nodeStore cache.Store,
	si step.Informer,
	cb step.UpdateStepCallback,
) *Controller {
	wq := workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "Step")
	controller := &Controller{
		schedulerName:  schedulerName,
		informer:       si,
		workqueue:      wq,
		workerNums:     4,
		cb:             cb,
		webhook:        hooks.NewDefaultStepWebHookPusher(),
		log:            zap.L().Named("Step"),
		runningWorkers: make(map[string]bool, 4),
	}

	si.Watcher().AddStepEventHandler(step.StepEventHandlerFuncs{
		AddFunc:    controller.enqueueForAdd,
		UpdateFunc: controller.enqueueForUpdate,
		DeleteFunc: controller.enqueueForDelete,
	})

	picker, err := roundrobin.NewStepPicker(nodeStore)
	if err != nil {
		panic(err)
	}
	controller.picker = picker
	return controller
}

// PipelineTaskScheduler 调度器控制器
type Controller struct {
	// workqueue is a rate limited work queue. This is used to queue work to be
	// processed instead of performing it as soon as a change happens. This
	// means we can ensure we only process a fixed amount of resources at a
	// time, and makes it easy to ensure we are never processing the same item
	// simultaneously in two different workers.
	workqueue      workqueue.RateLimitingInterface
	informer       step.Informer
	log            logger.Logger
	workerNums     int
	runningWorkers map[string]bool
	wLock          sync.Mutex
	picker         algorithm.StepPicker
	cb             step.UpdateStepCallback
	webhook        hooks.StepWebHookPusher
	schedulerName  string
}

func (c *Controller) SetWebHookPusher(p hooks.StepWebHookPusher) {
	c.webhook = p
}

// SetPicker 设置Node挑选器
func (c *Controller) SetStepPicker(picker algorithm.StepPicker) {
	c.picker = picker
}

func (c *Controller) Debug(log logger.Logger) {
	c.log = log
}

// Run will set up the event handlers for types we are interested in, as well
// as syncing informer caches and starting workers. It will block until stopCh
// is closed, at which point it will shutdown the workqueue and wait for
// workers to finish processing their current work items.
func (c *Controller) Run(ctx context.Context) error {
	// Start the informer factories to begin populating the informer caches
	c.log.Infof("starting step control loop, schedule name: %s", c.schedulerName)

	if err := c.sync(ctx); err != nil {
		return err
	}

	// 启动worker 处理来自Informer的事件
	for i := 0; i < c.workerNums; i++ {
		go c.runWorker(fmt.Sprintf("worker-%d", i))
	}

	c.waitDown(ctx)
	return nil
}

func (c *Controller) waitDown(ctx context.Context) {
	<-ctx.Done()
	// 关闭队列
	c.workqueue.ShutDown()
	// 停止worker
	c.log.Infof("step controller stopping, waitting for worker stop...")
	// 等待worker退出
	var max int
	for {
		if len(c.runningWorkers) == 0 {
			break
		}
		c.log.Infof("waiting worker %s exit ...", c.runningWorkerNames())
		if max > 30 {
			c.log.Warnf("waiting worker %s max times(30s) force exit", c.runningWorkerNames())
			break
		}
		max++
		time.Sleep(1 * time.Second)
	}
	c.log.Infof("step controller worker stopped commplet, now workers: %s", c.runningWorkerNames())
}

func (c *Controller) sync(ctx context.Context) error {
	// 获取所有的pipeline
	listCount := 0
	steps, err := c.informer.Lister().List(ctx)
	if err != nil {
		return err
	}

	// 看看是否有需要调度的
	for i := range steps {
		s := steps[i]
		if s.IsComplete() {
			c.log.Debugf("step %s is complete, skip schedule", s.Key)
			continue
		}

		if s.IsScheduled() {
			c.log.Debugf("step %s is scheduler to %s, skip schedule", s.Key, s.ScheduledNodeName())
			continue
		}
		c.informer.GetStore().Add(s)
		c.enqueueForAdd(s)
		listCount++
	}
	c.log.Infof("%d step need schedule", listCount)
	return nil
}

// runWorker is a long-running function that will continually call the
// processNextWorkItem function in order to read and process a message on the
// workqueue.
func (c *Controller) runWorker(name string) {
	isRunning, ok := c.runningWorkers[name]
	if ok && isRunning {
		c.log.Warnf("worker %s has running", name)
		return
	}
	c.wLock.Lock()
	c.runningWorkers[name] = true
	c.log.Infof("start worker %s", name)
	c.wLock.Unlock()
	for c.processNextWorkItem() {
	}
	if isRunning, ok = c.runningWorkers[name]; ok {
		c.wLock.Lock()
		delete(c.runningWorkers, name)
		c.wLock.Unlock()
		c.log.Infof("worker %s has stopped", name)
	}
}

// processNextWorkItem will read a single work item off the workqueue and
// attempt to process it, by calling the syncHandler.
func (c *Controller) processNextWorkItem() bool {
	obj, shutdown := c.workqueue.Get()
	if shutdown {
		return false
	}
	c.log.Debugf("get obj from queue: %s", obj)
	// We wrap this block in a func so we can defer c.workqueue.Done.
	err := func(obj interface{}) error {
		// We call Done here so the workqueue knows we have finished
		// processing this item. We also must remember to call Forget if we
		// do not want this work item being re-queued. For example, we do
		// not call Forget if a transient error occurs, instead the item is
		// put back on the workqueue and attempted again after a back-off
		// period.
		defer c.workqueue.Done(obj)
		var key string
		var ok bool

		// We expect strings to come off the workqueue. These are of the
		// form namespace/name. We do this as the delayed nature of the
		// workqueue means the items in the informer cache may actually be
		// more up to date that when the item was initially put onto the
		// workqueue.
		if key, ok = obj.(string); !ok {
			// As the item in the workqueue is actually invalid, we call
			// Forget here else we'd go into a loop of attempting to
			// process a work item that is invalid.
			c.workqueue.Forget(obj)
			c.log.Errorf("expected string in workqueue but got %#v", obj)
			return nil
		}
		c.log.Debugf("wait sync: %s", key)

		// Run the syncHandler, passing it the namespace/name string of the
		// Network resource to be synced.
		if err := c.syncHandler(key); err != nil {
			return fmt.Errorf("error syncing '%s': %s", key, err.Error())
		}
		// Finally, if no error occurs we Forget this item so it does not
		// get queued again until another change happens.
		c.workqueue.Forget(obj)
		c.log.Infof("successfully synced '%s'", key)
		return nil
	}(obj)
	if err != nil {
		c.log.Error(err)
		return true
	}
	return true
}

func (c *Controller) runningWorkerNames() string {
	c.wLock.Lock()
	defer c.wLock.Unlock()
	kList := make([]string, 0, len(c.runningWorkers))
	for k := range c.runningWorkers {
		kList = append(kList, k)
	}
	return strings.Join(kList, ",")
}

// enqueueNetwork takes a Cronjob resource and converts it into a namespace/name
// string which is then put onto the work queue. This method should *not* be
// passed resources of any type other than Cronjob.
func (c *Controller) enqueueForAdd(s *pipeline.Step) {
	c.log.Infof("receive add object: %s", s)
	if err := s.Validate(); err != nil {
		c.log.Errorf("invalidate *pipeline.Step obj")
		return
	}

	// 判断入队条件, 已经执行完的无需重复处理
	if s.IsComplete() {
		c.log.Errorf("step %s is complete, skip enqueue", s.Key)
		return
	}

	key := s.MakeObjectKey()
	c.workqueue.AddRateLimited(key)
}

// enqueueNetworkForDelete takes a deleted Network resource and converts it into a namespace/name
// string which is then put onto the work queue. This method should *not* be
// passed resources of any type other than Network.
func (c *Controller) enqueueForDelete(s *pipeline.Step) {
	c.log.Infof("receive delete object: %s", s)
	if err := s.Validate(); err != nil {
		c.log.Errorf("invalidate *pipeline.Step obj")
		return
	}
	key := s.MakeObjectKey()
	c.workqueue.AddRateLimited(key)
}

// 如果step有状态更新, 回调通知pipeline controller
func (c *Controller) enqueueForUpdate(oldObj, newObj *pipeline.Step) {
	c.log.Debugf("enqueue update old[%d], new[%d] ...", oldObj.ResourceVersion, newObj.ResourceVersion)

	// 判断事件状态, 调用webhook
	if err := c.webhook.Send(context.Background(), newObj.MatchedHooks(), newObj); err != nil {
		c.log.Errorf("send web hook error, %s", err)
	}

	switch newObj.CreateType {
	case pipeline.STEP_CREATE_BY_PIPELINE:
		// 如果是pipeline创建的，将事件传递给pipeline
		if c.cb != nil {
			c.cb(oldObj, newObj)
		}
	}

	key := newObj.MakeObjectKey()
	c.workqueue.AddRateLimited(key)
}

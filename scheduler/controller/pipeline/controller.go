package pipeline

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/logger/zap"
	"k8s.io/client-go/util/workqueue"

	"github.com/infraboard/workflow/api/apps/pipeline"
	"github.com/infraboard/workflow/common/cache"
	"github.com/infraboard/workflow/scheduler/algorithm"
	"github.com/infraboard/workflow/scheduler/algorithm/roundrobin"

	informer "github.com/infraboard/workflow/common/informers/pipeline"
	"github.com/infraboard/workflow/common/informers/step"
)

// NewPipelineController pipeline controller
func NewPipelineController(
	schedulerName string,
	nodeStore cache.Store,
	pi informer.Informer,
	si step.Informer,

) *Controller {
	wq := workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "Pipeline")
	controller := &Controller{
		schedulerName:  schedulerName,
		informer:       pi,
		step:           si,
		workqueue:      wq,
		workerNums:     4,
		log:            zap.L().Named("Pipeline"),
		runningWorkers: make(map[string]struct{}, 4),
	}

	pi.Watcher().AddPipelineTaskEventHandler(informer.PipelineTaskEventHandlerFuncs{
		AddFunc:    controller.enqueueForAdd,
		UpdateFunc: controller.enqueueForUpdate,
		DeleteFunc: controller.handleDelete,
	})

	picker, err := roundrobin.NewPipelinePicker(nodeStore)
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
	informer       informer.Informer
	step           step.Informer
	log            logger.Logger
	workerNums     int
	runningWorkers map[string]struct{}
	wLock          sync.Mutex
	picker         algorithm.PipelinePicker
	schedulerName  string
}

// SetPicker 设置Node挑选器
func (c *Controller) SetPipelinePicker(picker algorithm.PipelinePicker) {
	c.picker = picker
}

func (c *Controller) Debug(log logger.Logger) {
	c.log = log
}

func (c *Controller) Run(ctx context.Context) error {
	return c.run(ctx, false)
}

func (c *Controller) AsyncRun(ctx context.Context) error {
	return c.run(ctx, true)
}

// Run will set up the event handlers for types we are interested in, as well
// as syncing informer caches and starting workers. It will block until stopCh
// is closed, at which point it will shutdown the workqueue and wait for
// workers to finish processing their current work items.
func (c *Controller) run(ctx context.Context, async bool) error {
	// Start the informer factories to begin populating the informer caches
	c.log.Infof("starting pipeline control loop, schedule name: %s", c.schedulerName)

	// 获取所有的pipeline
	if err := c.sync(ctx); err != nil {
		return err
	}

	// 启动worker 处理来自Informer的事件
	for i := 0; i < c.workerNums; i++ {
		go c.runWorker(fmt.Sprintf("worker-%d", i))
	}

	if async {
		go c.waitDown(ctx)
	} else {
		c.waitDown(ctx)
	}

	return nil
}

func (c *Controller) waitDown(ctx context.Context) {
	<-ctx.Done()
	// 关闭队列
	c.workqueue.ShutDown()
	// 停止worker
	c.log.Infof("pipeline controller stopping, waitting for worker stop...")
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
	c.log.Infof("pipeline controller worker stopped commplet, now workers: %s", c.runningWorkerNames())
}

func (c *Controller) sync(ctx context.Context) error {
	// 获取所有的pipeline
	listCount := 0
	ps, err := c.informer.Lister().List(ctx, nil)
	if err != nil {
		return err
	}

	// 看看是否有需要调度的
	for i := range ps.Items {
		p := ps.Items[i]

		if p.IsComplete() {
			c.log.Debugf("pipline %s is complete, skip schedule",
				p.ShortDescribe())
			continue
		}

		if !p.MatchScheduler(c.schedulerName) {
			c.log.Debugf("pipeline %s scheduler %s is not match this scheduler %s",
				p.ShortDescribe(), p.ScheduledNodeName(), c.schedulerName)
			continue
		}

		c.enqueueForAdd(p)
		listCount++
	}
	c.log.Infof("%d pipeline need schedule", listCount)
	return nil
}

// runWorker is a long-running function that will continually call the
// processNextWorkItem function in order to read and process a message on the
// workqueue.
func (c *Controller) runWorker(name string) {
	_, ok := c.runningWorkers[name]
	if ok {
		c.log.Warnf("worker %s has running", name)
		return
	}
	c.wLock.Lock()
	c.runningWorkers[name] = struct{}{}
	c.log.Infof("start worker %s", name)
	c.wLock.Unlock()
	for c.processNextWorkItem() {
	}
	if _, ok = c.runningWorkers[name]; ok {
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
func (c *Controller) enqueueForAdd(p *pipeline.Pipeline) {
	c.log.Infof("receive add object: %s", p)
	if err := p.Validate(); err != nil {
		c.log.Errorf("validate pipeline error, %s", err)
		return
	}

	// 判断入队条件, 已经执行完的无需重复处理
	if p.IsComplete() {
		c.log.Errorf("pipeline %s is complete, skip enqueue", p.ShortDescribe())
		return
	}

	key := p.MakeObjectKey()
	c.workqueue.AddRateLimited(key)
}

// enqueueNetworkForDelete takes a deleted Network resource and converts it into a namespace/name
// string which is then put onto the work queue. This method should *not* be
// passed resources of any type other than Network.
func (c *Controller) handleDelete(p *pipeline.Pipeline) {
	c.log.Infof("receive delete object: %s", p)
	if err := p.Validate(); err != nil {
		c.log.Errorf("validate pipeline error, %s", err)
		return
	}

	fmt.Println(p)
}

// pipeline有状态变化, 并且状态变化是由step变化为comeplete引用的 则需要进行Next任务调度了
func (c *Controller) enqueueForUpdate(old, new *pipeline.Pipeline) {
	c.log.Infof("receive update object: old: %s, new: %s", old.ShortDescribe(), new.ShortDescribe())

	// 已经处理完成的无需处理
	if new.IsComplete() {
		c.log.Debugf("skip run complete pipeline %s, status: %s", new.ShortDescribe(), new.Status.Status)
		return
	}

	key := new.MakeObjectKey()
	c.workqueue.AddRateLimited(key)
}

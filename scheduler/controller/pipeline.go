package controller

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/infraboard/mcube/logger"
	"k8s.io/client-go/util/workqueue"

	"github.com/infraboard/workflow/api/pkg/pipeline"
	"github.com/infraboard/workflow/scheduler/algorithm"
	"github.com/infraboard/workflow/scheduler/algorithm/roundrobin"
	"github.com/infraboard/workflow/scheduler/informer"
	"github.com/infraboard/workflow/scheduler/store"
)

// PipelineScheduler pipeline controller
func NewPipelineScheduler(
	nodeStore store.NodeStore,
	inform informer.Informer,
) *PipelineScheduler {
	wq := workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "PipelineScheduler")
	controller := &PipelineScheduler{
		nodes:          nodeStore,
		workqueue:      wq,
		lister:         inform.Lister(),
		workerNums:     4,
		runningWorkers: make(map[string]bool, 4),
	}

	inform.Watcher().AddNodeEventHandler(informer.NodeEventHandlerFuncs{
		AddFunc:    controller.addNode,
		DeleteFunc: controller.delNode,
	})
	inform.Watcher().AddPipelineEventHandler(informer.PipelineEventHandlerFuncs{
		AddFunc: controller.addPipeline,
	})
	inform.Watcher().AddStepEventHandler(informer.StepEventHandlerFuncs{
		AddFunc: controller.addStep,
	})

	picker, err := roundrobin.NewPicker(nodeStore)
	if err != nil {
		panic(err)
	}
	controller.picker = picker
	return controller
}

// SchedulerController 调度器控制器
type PipelineScheduler struct {
	// workqueue is a rate limited work queue. This is used to queue work to be
	// processed instead of performing it as soon as a change happens. This
	// means we can ensure we only process a fixed amount of resources at a
	// time, and makes it easy to ensure we are never processing the same item
	// simultaneously in two different workers.
	workqueue      workqueue.RateLimitingInterface
	lister         informer.Lister
	log            logger.Logger
	workerNums     int
	runningWorkers map[string]bool
	wLock          sync.Mutex
	picker         algorithm.Picker
	nodes          store.NodeStore // 存储每个region的node信息
}

// SetPicker 设置Node挑选器
func (c *PipelineScheduler) SetPicker(picker algorithm.Picker) {
	c.picker = picker
}

// Run will set up the event handlers for types we are interested in, as well
// as syncing informer caches and starting workers. It will block until stopCh
// is closed, at which point it will shutdown the workqueue and wait for
// workers to finish processing their current work items.
func (c *PipelineScheduler) Run(ctx context.Context) error {
	// Start the informer factories to begin populating the informer caches
	c.log.Info("Starting schedule control loop")
	// 调用Lister 获得所有的cronjob 并添加cron
	c.log.Info("Starting Sync(List) All nodes")
	nodes, err := c.lister.ListNode(ctx)
	if err != nil {
		return err
	}
	c.nodes.AddNodeSet(nodes)
	c.log.Infof("Sync All(%d) nodes success", len(nodes))
	// 获取所有的pipeline
	listCount := 0
	ps, err := c.lister.ListPipeline(ctx, nil)
	if err != nil {
		return err
	}

	// 看看是否有需要调度的
	for i := range ps.Items {
		if ps.Items[i].ScheduledNodeName() == "" {
			c.workqueue.Add(ps.Items[i])
			listCount++
		}
	}
	c.log.Infof("%d pipeline need scheduled", listCount)

	// 启动worker 处理来自Informer的事件
	for i := 0; i < c.workerNums; i++ {
		go c.runWorker(fmt.Sprintf("worker-%d", i))
	}
	<-ctx.Done()
	// 关闭队列
	c.workqueue.ShutDown()
	// 停止worker
	c.log.Infof("scheduler controller stopping, waitting for worker stop...")
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
	c.log.Infof("scheduler controller worker stopped commplet, now workers: %s", c.runningWorkerNames())
	return nil
}

// runWorker is a long-running function that will continually call the
// processNextWorkItem function in order to read and process a message on the
// workqueue.
func (c *PipelineScheduler) runWorker(name string) {
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
	return
}

// processNextWorkItem will read a single work item off the workqueue and
// attempt to process it, by calling the syncHandler.
func (c *PipelineScheduler) processNextWorkItem() bool {
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
		switch v := obj.(type) {
		case *pipeline.Pipeline:
			c.log.Debugf("wait schedule pipeline: %s", v.GetId())
			if err := c.schedulePipeline(v); err != nil {
				return fmt.Errorf("error scheduled '%s': %s", v, err.Error())
			}
		case *pipeline.Step:
			c.log.Debugf("wait schedule step: %s", v.GetId())
			if err := c.scheduleStep(v); err != nil {
				return fmt.Errorf("error scheduled '%s': %s", v, err.Error())
			}
		default:
			// As the item in the workqueue is actually invalid, we call
			// Forget here else we'd go into a loop of attempting to
			// process a work item that is invalid.
			c.workqueue.Forget(obj)
			c.log.Errorf("unknow type for schedule: %s in workqueue", reflect.TypeOf(obj))
			return nil
		}

		// Finally, if no error occurs we Forget this item so it does not
		// get queued again until another change happens.
		c.workqueue.Forget(obj)
		c.log.Infof("Successfully scheduled '%s'", obj)
		return nil
	}(obj)
	if err != nil {
		c.log.Error(err)
		return true
	}
	return true
}

func (c *PipelineScheduler) runningWorkerNames() string {
	c.wLock.Lock()
	defer c.wLock.Unlock()
	kList := make([]string, 0, len(c.runningWorkers))
	for k := range c.runningWorkers {
		kList = append(kList, k)
	}
	return strings.Join(kList, ",")
}

func (c *PipelineScheduler) schedulePipeline(pp *pipeline.Pipeline) error {
	return nil
}

func (c *PipelineScheduler) scheduleStep(step *pipeline.Step) error {
	node, err := c.picker.Pick(step)
	if err != nil {
		return err
	}

	// 没有合法的node
	if node == nil {
		return fmt.Errorf("no excutable node")
	}

	c.log.Debugf("choice node %s for step %s", node.InstanceName, step.Id)
	step.AddScheduleNode(node.InstanceName)
	// 清除一下其他数据
	if err := c.lister.UpdateStep(step); err != nil {
		c.log.Errorf("update scheduled step error, %s", err)
	}
	return nil
}

package node

import (
	"context"

	"github.com/infraboard/workflow/api/pkg/node"
	"github.com/infraboard/workflow/common/cache"
)

// CronJobInformer provides access to a shared informer and lister for
// CronJobs.
type Informer interface {
	Watcher() Watcher
	// List All Node
	Lister() Lister
	// 获取缓存
	GetStore() cache.Store
}

// Lister 获取所有执行节点
type Lister interface {
	// List lists all Node
	List(context.Context) ([]*node.Node, error)
}

type Watcher interface {
	// Run starts and runs the shared informer, returning after it stops.
	// The informer will be stopped when stopCh is closed.
	Run(ctx context.Context) error
	// AddEventHandler adds an event handler to the shared informer using the shared informer's resync
	// period.  Events to a single handler are delivered sequentially, but there is no coordination
	// between different handlers.
	AddNodeEventHandler(handler NodeEventHandler)
}

// NodeEventHandler can handle notifications for events that happen to a
// resource. The events are informational only, so you can't return an
// error.
//  * OnAdd is called when an object is added.
//  * OnUpdate is called when an object is modified. Note that oldObj is the
//      last known state of the object-- it is possible that several changes
//      were combined together, so you can't use this to see every single
//      change. OnUpdate is also called when a re-list happens, and it will
//      get called even if nothing changed. This is useful for periodically
//      evaluating or syncing something.
//  * OnDelete will get the final state of the item if it is known, otherwise
//      it will get an object of type DeletedFinalStateUnknown. This can
//      happen if the watch is closed and misses the delete event and we don't
//      notice the deletion until the subsequent re-list.
type NodeEventHandler interface {
	OnAdd(node *node.Node)
	OnUpdate(oldNode, newNode *node.Node)
	OnDelete(node *node.Node)
}

// NodeEventHandlerFuncs is an adaptor to let you easily specify as many or
// as few of the notification functions as you want while still implementing
// ResourceEventHandler.
type NodeEventHandlerFuncs struct {
	AddFunc    func(obj *node.Node)
	UpdateFunc func(oldObj, newObj *node.Node)
	DeleteFunc func(obj *node.Node)
}

// OnAdd calls AddFunc if it's not nil.
func (r NodeEventHandlerFuncs) OnAdd(obj *node.Node) {
	if r.AddFunc != nil {
		r.AddFunc(obj)
	}
}

// OnUpdate calls UpdateFunc if it's not nil.
func (r NodeEventHandlerFuncs) OnUpdate(oldObj, newObj *node.Node) {
	if r.UpdateFunc != nil {
		r.UpdateFunc(oldObj, newObj)
	}
}

// OnDelete calls DeleteFunc if it's not nil.
func (r NodeEventHandlerFuncs) OnDelete(obj *node.Node) {
	if r.DeleteFunc != nil {
		r.DeleteFunc(obj)
	}
}

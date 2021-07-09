package step

import (
	"context"

	"github.com/infraboard/workflow/api/pkg/pipeline"
	"github.com/infraboard/workflow/common/cache"
)

// Informer 负责事件通知
type Informer interface {
	Watcher() Watcher
	Lister() Lister
	Recorder() Recorder
	GetStore() cache.Store
}

type Recorder interface {
	Update(*pipeline.Step) error
}

type Lister interface {
	List(ctx context.Context) ([]*pipeline.Step, error)
}

// Watcher 负责事件通知
type Watcher interface {
	// Run starts and runs the shared informer, returning after it stops.
	// The informer will be stopped when stopCh is closed.
	Run(ctx context.Context) error
	// AddEventHandler adds an event handler to the shared informer using the shared informer's resync
	// period.  Events to a single handler are delivered sequentially, but there is no coordination
	// between different handlers.
	AddStepEventHandler(handler StepEventHandler)
	AddStepFilterHandler(handler StepFilterHandler)
}

// StepEventHandler can handle notifications for events that happen to a
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
type StepEventHandler interface {
	OnAdd(obj *pipeline.Step)
	OnUpdate(old, new *pipeline.Step)
	OnDelete(obj *pipeline.Step)
}

// StepEventHandlerFuncs is an adaptor to let you easily specify as many or
// as few of the notification functions as you want while still implementing
// ResourceEventHandler.
type StepEventHandlerFuncs struct {
	AddFunc    func(obj *pipeline.Step)
	UpdateFunc func(oldObj, newObj *pipeline.Step)
	DeleteFunc func(obj *pipeline.Step)
}

// OnAdd calls AddFunc if it's not nil.
func (r StepEventHandlerFuncs) OnAdd(obj *pipeline.Step) {
	if r.AddFunc != nil {
		r.AddFunc(obj)
	}
}

// OnUpdate calls UpdateFunc if it's not nil.
func (r StepEventHandlerFuncs) OnUpdate(oldObj, newObj *pipeline.Step) {
	if r.UpdateFunc != nil {
		r.UpdateFunc(oldObj, newObj)
	}
}

// OnDelete calls DeleteFunc if it's not nil.
func (r StepEventHandlerFuncs) OnDelete(obj *pipeline.Step) {
	if r.DeleteFunc != nil {
		r.DeleteFunc(obj)
	}
}

type StepFilterHandler func(obj *pipeline.Step) bool

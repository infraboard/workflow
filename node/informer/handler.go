package informer

import (
	"github.com/infraboard/workflow/api/pkg/node"
	"github.com/infraboard/workflow/api/pkg/pipeline"
)

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

// PipelineEventHandler can handle notifications for events that happen to a
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
type PipelineTaskEventHandler interface {
	OnAdd(obj *pipeline.Pipeline)
	OnUpdate(old, new *pipeline.Pipeline)
	OnDelete(obj *pipeline.Pipeline)
}

// PipelineEventHandlerFuncs is an adaptor to let you easily specify as many or
// as few of the notification functions as you want while still implementing
// ResourceEventHandler.
type PipelineTaskEventHandlerFuncs struct {
	AddFunc    func(obj *pipeline.Pipeline)
	UpdateFunc func(oldObj, newObj *pipeline.Pipeline)
	DeleteFunc func(obj *pipeline.Pipeline)
}

// OnAdd calls AddFunc if it's not nil.
func (r PipelineTaskEventHandlerFuncs) OnAdd(obj *pipeline.Pipeline) {
	if r.AddFunc != nil {
		r.AddFunc(obj)
	}
}

// OnUpdate calls UpdateFunc if it's not nil.
func (r PipelineTaskEventHandlerFuncs) OnUpdate(oldObj, newObj *pipeline.Pipeline) {
	if r.UpdateFunc != nil {
		r.UpdateFunc(oldObj, newObj)
	}
}

// OnDelete calls DeleteFunc if it's not nil.
func (r PipelineTaskEventHandlerFuncs) OnDelete(obj *pipeline.Pipeline) {
	if r.DeleteFunc != nil {
		r.DeleteFunc(obj)
	}
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

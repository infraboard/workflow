package informer

import "context"

// Watcher 负责事件通知
type Watcher interface {
	// Run starts and runs the shared informer, returning after it stops.
	// The informer will be stopped when stopCh is closed.
	Run(ctx context.Context) error
	// AddEventHandler adds an event handler to the shared informer using the shared informer's resync
	// period.  Events to a single handler are delivered sequentially, but there is no coordination
	// between different handlers.
	AddNodeEventHandler(handler NodeEventHandler)
	AddPipelineTaskEventHandler(handler PipelineTaskEventHandler)
	AddStepEventHandler(handler StepEventHandler)
}

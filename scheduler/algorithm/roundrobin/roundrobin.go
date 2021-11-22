package roundrobin

import (
	"fmt"
	"sync"

	"github.com/infraboard/workflow/api/app/node"
	"github.com/infraboard/workflow/api/app/pipeline"
	"github.com/infraboard/workflow/common/cache"
	"github.com/infraboard/workflow/scheduler/algorithm"
)

type roundrobinPicker struct {
	mu    *sync.Mutex
	next  int
	store cache.Store
}

// NewStepPicker 实现分调度
func NewStepPicker(nodestore cache.Store) (algorithm.StepPicker, error) {
	base := &roundrobinPicker{
		store: nodestore,
		mu:    new(sync.Mutex),
		next:  0,
	}
	return &stepPicker{base}, nil
}

type stepPicker struct {
	*roundrobinPicker
}

func (p *stepPicker) Pick(step *pipeline.Step) (*node.Node, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	nodes := p.store.List()
	if len(nodes) == 0 {
		return nil, fmt.Errorf("has no available nodes")
	}

	n := nodes[p.next]

	// 修改状态
	p.next = (p.next + 1) % p.store.Len()

	return n.(*node.Node), nil
}

// NewPipelinePicker 实现分调度
func NewPipelinePicker(nodestore cache.Store) (algorithm.PipelinePicker, error) {
	base := &roundrobinPicker{
		store: nodestore,
		mu:    new(sync.Mutex),
		next:  0,
	}
	return &pipelinePicker{base}, nil
}

type pipelinePicker struct {
	*roundrobinPicker
}

func (p *pipelinePicker) Pick(step *pipeline.Pipeline) (*node.Node, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	nodes := p.store.List()
	if len(nodes) == 0 {
		return nil, fmt.Errorf("has no available nodes")
	}

	n := nodes[p.next]

	// 修改状态
	p.next = (p.next + 1) % p.store.Len()

	return n.(*node.Node), nil
}

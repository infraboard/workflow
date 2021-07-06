package roundrobin

import (
	"fmt"
	"sync"

	"github.com/infraboard/workflow/api/pkg/node"
	"github.com/infraboard/workflow/api/pkg/pipeline"
	"github.com/infraboard/workflow/scheduler/algorithm"
	"github.com/infraboard/workflow/scheduler/store"
)

type roundrobinPicker struct {
	mu    *sync.Mutex
	next  int
	store store.NodeStore
}

// NewStepPicker 实现分调度
func NewStepPicker(nodestore store.NodeStore) (algorithm.StepPicker, error) {
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

	nodes := p.store.ListNode()
	if len(nodes) == 0 {
		return nil, fmt.Errorf("has no available nodes")
	}

	node := nodes[p.next]

	// 修改状态
	p.next = (p.next + 1) % p.store.Len()

	return node, nil
}

// NewPipelinePicker 实现分调度
func NewPipelinePicker(nodestore store.NodeStore) (algorithm.PipelinePicker, error) {
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

	nodes := p.store.ListNode()
	if len(nodes) == 0 {
		return nil, fmt.Errorf("has no available nodes")
	}

	node := nodes[p.next]

	// 修改状态
	p.next = (p.next + 1) % p.store.Len()

	return node, nil
}

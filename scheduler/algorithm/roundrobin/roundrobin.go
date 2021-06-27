package roundrobin

import (
	"errors"
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
	if p.store.Len() == 0 {
		return nil, errors.New("no available node")
	}

	if p.store.Len() == 0 {
		return nil, fmt.Errorf("has no available nodes")
	}
	p.mu.Lock()
	defer p.mu.Unlock()

	node, err := p.store.GetByIndex(p.next)
	if err != nil {
		return nil, err
	}

	// 修改状态
	p.next = (p.next + 1) % p.store.Len()

	return node, nil
}

// NewTaskPicker 实现分调度
func NewTaskPicker(nodestore store.NodeStore) (algorithm.TaskPicker, error) {
	base := &roundrobinPicker{
		store: nodestore,
		mu:    new(sync.Mutex),
		next:  0,
	}
	return &taskPicker{base}, nil
}

type taskPicker struct {
	*roundrobinPicker
}

func (p *taskPicker) Pick(step *pipeline.Pipeline) (*node.Node, error) {
	if p.store.Len() == 0 {
		return nil, errors.New("no available node")
	}

	if p.store.Len() == 0 {
		return nil, fmt.Errorf("has no available nodes")
	}
	p.mu.Lock()
	defer p.mu.Unlock()

	node, err := p.store.GetByIndex(p.next)
	if err != nil {
		return nil, err
	}

	// 修改状态
	p.next = (p.next + 1) % p.store.Len()

	return node, nil
}

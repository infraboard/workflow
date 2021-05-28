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

// NewPicker 实现分调度
func NewPicker(nodestore store.NodeStore) (algorithm.Picker, error) {
	return &roundrobinPicker{
		store: nodestore,
		mu:    new(sync.Mutex),
		next:  0,
	}, nil
}

type roundrobinPicker struct {
	mu    *sync.Mutex
	next  int
	store store.NodeStore
}

func (p *roundrobinPicker) Pick(step *pipeline.Step) (*node.Node, error) {
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

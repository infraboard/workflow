package store

import (
	"container/list"
	"fmt"
	"sync"

	"github.com/infraboard/workflow/api/pkg/node"
)

// NodeStore 保持node
type NodeStore interface {
	ListNode() []*node.Node
	AddNodeSet([]*node.Node)
	AddNode(*node.Node)
	DelNode(*node.Node) error
	GetByIndex(i int) (*node.Node, error)
	Len() int
}

// NewDeaultNodeStore 使用默认存储
func NewDeaultNodeStore() NodeStore {
	return &nodeStore{
		nodes: list.New(),
		mu:    new(sync.Mutex),
	}
}

// key 为region
type nodeStore struct {
	nodes *list.List
	mu    *sync.Mutex
}

func (s *nodeStore) Len() int {
	return s.nodes.Len()
}

func (s *nodeStore) ListNode() []*node.Node {
	nodes := make([]*node.Node, 0, s.Len())
	for e := s.nodes.Front(); e != nil; e = e.Next() {
		nodes = append(nodes, e.Value.(*node.Node))
	}

	return nodes
}

func (s *nodeStore) GetByIndex(i int) (*node.Node, error) {
	if i > s.Len()-1 {
		return nil, fmt.Errorf("out of index")
	}

	var index int
	for e := s.nodes.Front(); e != nil; e = e.Next() {
		if i == index {
			return e.Value.(*node.Node), nil
		}
		index++
	}

	return nil, fmt.Errorf("not found")
}

func (s *nodeStore) GetByName(name string) (*node.Node, error) {
	for e := s.nodes.Front(); e != nil; e = e.Next() {
		n := e.Value.(*node.Node)
		if n.Name() == name {
			return n, nil
		}
	}

	return nil, fmt.Errorf("%s node no found", name)
}

func (s *nodeStore) AddNodeSet(nodes []*node.Node) {
	for i := range nodes {
		s.AddNode(nodes[i])
	}
}

func (s *nodeStore) AddNode(n *node.Node) {
	if n != nil {
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.nodes.PushBack(n)
}

func (s *nodeStore) DelNode(n *node.Node) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.Len() == 0 {
		return fmt.Errorf("no node to delete")
	}

	for e := s.nodes.Back(); e != nil; e = e.Next() {
		if n.Name() == e.Value.(*node.Node).Name() {
			s.nodes.Remove(e)
			return nil
		}
	}

	return fmt.Errorf("not found")
}

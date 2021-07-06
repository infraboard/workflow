package store

import (
	"sync"

	"github.com/infraboard/workflow/api/pkg/node"
)

// NodeStore 保持node
type NodeStore interface {
	ListNode() []*node.Node
	AddNodeSet([]*node.Node)
	AddNode(*node.Node)
	DelNode(*node.Node)
	Len() int
	LenWithRegion(region string) int
}

// NewDeaultNodeStore 使用默认存储
func NewDeaultNodeStore() NodeStore {
	return &mapNodeStore{
		nodes: map[string][]*node.Node{},
		mu:    new(sync.Mutex),
	}
}

// key 为region
type mapNodeStore struct {
	nodes map[string][]*node.Node
	mu    *sync.Mutex
}

func (s *mapNodeStore) Len() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	var count int
	for key := range s.nodes {
		count += len(s.nodes[key])
	}
	return count
}

func (s *mapNodeStore) LenWithRegion(region string) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	v, ok := s.nodes[region]
	if !ok {
		return 0
	}
	return len(v)
}

func (s *mapNodeStore) ListNode() []*node.Node {
	s.mu.Lock()
	defer s.mu.Unlock()

	ns := []*node.Node{}

	for _, v := range s.nodes {
		ns = append(ns, v...)
	}

	return ns
}

func (s *mapNodeStore) AddNodeSet(nodes []*node.Node) {
	for i := range nodes {
		s.AddNode(nodes[i])
	}
}

func (s *mapNodeStore) AddNode(n *node.Node) {
	if n == nil {
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.nodes[n.Region]; !ok {
		s.nodes[n.Region] = []*node.Node{}
	}
	s.nodes[n.Region] = append(s.nodes[n.Region], n)
}

func (s *mapNodeStore) DelNode(node *node.Node) {
	s.mu.Lock()
	defer s.mu.Unlock()
	rn, ok := s.nodes[node.Region]
	if !ok {
		return
	}
	for i := 0; i < len(rn); i++ {
		if rn[i].Name() == node.Name() {
			rn = append(rn[:i], rn[i+1:]...)
			i--
		}
	}
	s.nodes[node.Region] = rn
}

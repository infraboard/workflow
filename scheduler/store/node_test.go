package store_test

import (
	"fmt"
	"testing"

	"github.com/infraboard/workflow/api/pkg/node"
	"github.com/infraboard/workflow/scheduler/store"
)

var (
	nodeStore store.NodeStore
)

func TestNodeList(t *testing.T) {
	fmt.Println(nodeStore.ListNode())
	node2 := &node.Node{
		InstanceName: "node01",
		ServiceName:  "workflow",
		Type:         node.NodeType,
	}
	nodeStore.AddNodeSet([]*node.Node{node2})
	fmt.Println(nodeStore.ListNode())
	nodeStore.DelNode(node2)
	fmt.Println(nodeStore.ListNode())
}
func init() {
	nodeStore = store.NewDeaultNodeStore()
	node1 := &node.Node{
		InstanceName: "node01",
		ServiceName:  "workflow",
		Type:         node.NodeType,
	}
	nodeStore.AddNodeSet([]*node.Node{
		node1,
	})
}

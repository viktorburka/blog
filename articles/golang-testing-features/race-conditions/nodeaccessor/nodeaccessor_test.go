package nodeaccessor

import (
	"fmt"
	"testing"
)

func TestUpdateNodeNonConcurrent(t *testing.T) {
	type TestData struct {
		nodeName          string
		expectedNodeCount int
	}
	tests := []TestData{
		{nodeName: "/stream/1", expectedNodeCount: 0},
		{nodeName: "/stream/2", expectedNodeCount: 0},
		{nodeName: "/stream/1", expectedNodeCount: 1},
		{nodeName: "/stream/1", expectedNodeCount: 2},
	}
	accessor := NewNodeAccessor()
	for _, test := range tests {
		accessor.UpdateNode(test.nodeName)
		nodeCount, err := accessor.NodeCount(test.nodeName)
		if err != nil {
			t.Fatalf("node %v not found", test.nodeName)
		}
		if nodeCount != test.expectedNodeCount {
			t.Fatalf("incorrect node count for node %v: expected %v but got %v",
				test.nodeName, test.expectedNodeCount, nodeCount)
		}
	}
}

func TestUpdateNodeConcurrent(t *testing.T) {
	type TestData struct {
		nodeName          string
		expectedNodeCount int
	}
	tests := []TestData{
		{nodeName: "/stream/1", expectedNodeCount: 0},
		{nodeName: "/stream/2", expectedNodeCount: 0},
		{nodeName: "/stream/1", expectedNodeCount: 1},
		{nodeName: "/stream/1", expectedNodeCount: 2},
	}
	accessor := NewNodeAccessor()
	for idx, test := range tests {
		t.Run(fmt.Sprintf("test#%v", idx), func(t *testing.T) {
			test := test // capture range variable
			t.Parallel()
			accessor.UpdateNode(test.nodeName)
			nodeCount, err := accessor.NodeCount(test.nodeName)
			if err != nil {
				t.Fatalf("node %v not found", test.nodeName)
			}
			if nodeCount != test.expectedNodeCount {
				t.Fatalf("incorrect node count for node %v: expected %v but got %v",
					test.nodeName, test.expectedNodeCount, nodeCount)
			}
		})
	}
}





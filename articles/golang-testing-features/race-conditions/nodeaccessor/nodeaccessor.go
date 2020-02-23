package nodeaccessor

import (
	"fmt"
	"sync"
)

type NodeAccessor struct {
	counter map[string]int
	mux     sync.Mutex
}

func NewNodeAccessor() *NodeAccessor {
	return &NodeAccessor{
		counter: make(map[string]int),
	}
}

func (na *NodeAccessor) UpdateNode(name string) {
	na.mux.Lock()
	defer na.mux.Unlock()

	if _, ok := na.counter[name]; !ok {
		na.counter[name] = 0
	} else {
		na.counter[name]++
	}
}

func (na *NodeAccessor) NodeCount(name string) (int, error) {
	na.mux.Lock()
	defer na.mux.Unlock()

	counter, ok := na.counter[name]
	if !ok {
		return -1, fmt.Errorf("node does not exist")
	}
	return counter, nil
}

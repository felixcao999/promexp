package prometheus

import (
	"sync"
)

type NodeMetrics struct {
	sync.RWMutex
	metrics map[string]map[string]string
}

func NewNodeMetrics() *NodeMetrics {
	return &NodeMetrics{
		metrics: map[string]map[string]string{},
	}
}

func (nm *NodeMetrics) Set(node_addr, metric, value string) {
	nm.Lock()
	defer nm.Unlock()
	node_vs, ok := nm.metrics[node_addr]
	if !ok {
		node_vs = make(map[string]string)
		nm.metrics[node_addr] = node_vs
	}
	node_vs[metric] = value
}

func (nm *NodeMetrics) Get(node_addr, metric string) (string, bool) {
	nm.RLock()
	defer nm.RUnlock()
	node_vs, ok := nm.metrics[node_addr]
	if ok {
		v, ok := node_vs[metric]
		if ok {
			return v, true
		}
	}
	return "", false
}

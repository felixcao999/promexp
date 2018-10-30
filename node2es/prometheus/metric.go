package prometheus

import (
	"math"
	"sync"
)

type NodeMetrics struct {
	sync.RWMutex
	metrics map[string]map[string]float64
}

func NewNodeMetrics() *NodeMetrics {
	return &NodeMetrics{
		metrics: map[string]map[string]float64{},
	}
}

func (nm *NodeMetrics) Set(node_addr, metric string, value float64) {
	nm.Lock()
	defer nm.Unlock()
	node_vs, ok := nm.metrics[node_addr]
	if !ok {
		node_vs = make(map[string]float64)
		nm.metrics[node_addr] = node_vs
	}
	node_vs[metric] = round(value, 2)
}

func round(f float64, n int) float64 {
	n10 := math.Pow10(n)
	return math.Trunc((f+0.5/n10)*n10) / n10
}

func (nm *NodeMetrics) Get(node_addr, metric string) (float64, bool) {
	nm.RLock()
	defer nm.RUnlock()
	node_vs, ok := nm.metrics[node_addr]
	if ok {
		v, ok := node_vs[metric]
		if ok {
			return v, true
		}
	}
	return 0, false
}

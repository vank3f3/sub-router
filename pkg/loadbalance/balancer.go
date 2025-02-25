package loadbalance

import (
	"sync"
	"sync/atomic"
)

// Strategy 负载均衡策略
type Strategy int

const (
	RoundRobin Strategy = iota // 轮询
	Random                     // 随机
	WeightedRR                 // 加权轮询
)

// Backend 后端服务器
type Backend struct {
	URL      string // 服务器地址
	Weight   int    // 权重
	Healthy  bool   // 健康状态
	Active   int64  // 活跃连接数
	Priority int    // 优先级
}

// Balancer 负载均衡器接口
type Balancer interface {
	Add(backend *Backend)
	Remove(url string)
	Next() *Backend
	MarkDown(url string)
	MarkUp(url string)
}

// RoundRobinBalancer 轮询负载均衡器
type RoundRobinBalancer struct {
	backends []*Backend
	current  uint64
	mu       sync.RWMutex
}

// NewRoundRobinBalancer 创建轮询负载均衡器
func NewRoundRobinBalancer() *RoundRobinBalancer {
	return &RoundRobinBalancer{
		backends: make([]*Backend, 0),
	}
}

// Add 添加后端服务器
func (b *RoundRobinBalancer) Add(backend *Backend) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.backends = append(b.backends, backend)
}

// Remove 移除后端服务器
func (b *RoundRobinBalancer) Remove(url string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	for i, backend := range b.backends {
		if backend.URL == url {
			b.backends = append(b.backends[:i], b.backends[i+1:]...)
			return
		}
	}
}

// Next 获取下一个后端服务器
func (b *RoundRobinBalancer) Next() *Backend {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if len(b.backends) == 0 {
		return nil
	}

	// 获取可用的后端服务器
	var availableBackends []*Backend
	for _, backend := range b.backends {
		if backend.Healthy {
			availableBackends = append(availableBackends, backend)
		}
	}

	if len(availableBackends) == 0 {
		return nil
	}

	// 原子操作获取下一个索引
	next := atomic.AddUint64(&b.current, 1)
	idx := next % uint64(len(availableBackends))
	return availableBackends[idx]
}

// MarkDown 标记服务器为不可用
func (b *RoundRobinBalancer) MarkDown(url string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	for _, backend := range b.backends {
		if backend.URL == url {
			backend.Healthy = false
			return
		}
	}
}

// MarkUp 标记服务器为可用
func (b *RoundRobinBalancer) MarkUp(url string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	for _, backend := range b.backends {
		if backend.URL == url {
			backend.Healthy = true
			return
		}
	}
}

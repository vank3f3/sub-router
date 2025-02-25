package benchmark

import (
	"net/http/httptest"
	"sync/atomic"
)

// TestBalancer 用于测试的简单负载均衡器
type TestBalancer struct {
	backends []*httptest.Server
	current  uint64
}

// NewTestBalancer 创建测试用负载均衡器
func NewTestBalancer(backends []*httptest.Server) *TestBalancer {
	return &TestBalancer{
		backends: backends,
	}
}

// Next 获取下一个后端服务器
func (b *TestBalancer) Next() *httptest.Server {
	next := atomic.AddUint64(&b.current, 1)
	return b.backends[next%uint64(len(b.backends))]
}

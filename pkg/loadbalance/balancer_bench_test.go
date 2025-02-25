package loadbalance

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func BenchmarkLoadBalancer(b *testing.B) {
	// 创建多个后端服务器
	backends := make([]*Backend, 3)
	servers := make([]*httptest.Server, 3)

	for i := range backends {
		servers[i] = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("test response"))
		}))
		defer servers[i].Close()

		backends[i] = &Backend{
			URL:     servers[i].URL,
			Weight:  1,
			Healthy: true,
		}
	}

	// 创建负载均衡器
	balancer := NewRoundRobinBalancer()
	for _, backend := range backends {
		balancer.Add(backend)
	}

	// 运行基准测试
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			backend := balancer.Next()
			resp, err := http.Get(backend.URL)
			if err == nil {
				resp.Body.Close()
			}
		}
	})
}

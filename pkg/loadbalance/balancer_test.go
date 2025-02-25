package loadbalance

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRoundRobinBalancer(t *testing.T) {
	// 创建测试服务器
	server1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("server1"))
	}))
	defer server1.Close()

	server2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("server2"))
	}))
	defer server2.Close()

	// 创建负载均衡器
	balancer := NewRoundRobinBalancer()

	// 添加后端服务器
	balancer.Add(&Backend{
		URL:     server1.URL,
		Weight:  1,
		Healthy: true,
	})
	balancer.Add(&Backend{
		URL:     server2.URL,
		Weight:  1,
		Healthy: true,
	})

	// 测试轮询分发
	backends := make(map[string]int)
	for i := 0; i < 10; i++ {
		backend := balancer.Next()
		backends[backend.URL]++
	}

	// 验证分发是否均匀
	if backends[server1.URL] != 5 || backends[server2.URL] != 5 {
		t.Errorf("Load balancing is not even: %v", backends)
	}

	// 测试服务器标记为不可用
	balancer.MarkDown(server1.URL)
	backend := balancer.Next()
	if backend.URL != server2.URL {
		t.Errorf("Expected server2 after server1 marked down, got %s", backend.URL)
	}

	// 测试移除服务器
	balancer.Remove(server2.URL)
	if backend := balancer.Next(); backend != nil {
		t.Error("Expected nil after all servers removed")
	}
}

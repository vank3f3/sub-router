package benchmark

import (
	"net/http"
	"net/http/httptest"
	"sub-router/internal/config"
	"testing"

	"sub-router/internal/handler"

	"github.com/gin-gonic/gin"
)

func BenchmarkProxyHandler(b *testing.B) {
	// 设置测试配置
	config.GlobalConfig = config.Config{
		APIMappings: map[string]string{
			"test": "http://localhost:8888",
		},
	}

	// 创建测试服务器
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("test response"))
	}))
	defer backend.Close()

	// 更新测试配置
	config.GlobalConfig.APIMappings["test"] = backend.URL

	// 创建路由
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Any("/:service/*path", handler.ProxyHandler)

	// 运行基准测试
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/test/api/data", nil)
			router.ServeHTTP(w, req)
		}
	})
}

func BenchmarkLoadBalancer(b *testing.B) {
	// 创建多个后端服务器
	backends := make([]*httptest.Server, 3)
	for i := range backends {
		backends[i] = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("test response"))
		}))
		defer backends[i].Close()
	}

	// 设置负载均衡器
	balancer := NewTestBalancer(backends)

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

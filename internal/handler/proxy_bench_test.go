package handler

import (
	"net/http"
	"net/http/httptest"
	"sub-router/internal/config"
	"testing"

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
	router.Any("/:service/*path", ProxyHandler)

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

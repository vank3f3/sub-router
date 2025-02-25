package integration

import (
	"net/http"
	"net/http/httptest"
	"sub-router/internal/config"
	"testing"
	"time"

	"sub-router/pkg/loadbalance"
)

// setupTestServer 创建测试服务器
func setupTestServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/test/api":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		case "/test/slow":
			time.Sleep(200 * time.Millisecond)
			w.WriteHeader(http.StatusOK)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
}

// executeTest 执行测试请求
func executeTest(path string, timeout time.Duration) (*http.Response, error) {
	client := &http.Client{
		Timeout: timeout,
	}
	return client.Get("http://localhost:8080" + path)
}

func TestProxyIntegration(t *testing.T) {
	// 启动测试服务器
	testServer := setupTestServer()
	defer testServer.Close()

	// 配置代理
	config.GlobalConfig = config.Config{
		APIMappings: map[string]string{
			"test": testServer.URL, // 使用测试服务器的 URL
		},
	}

	// 创建负载均衡器
	balancer := loadbalance.NewRoundRobinBalancer()
	balancer.Add(&loadbalance.Backend{
		URL:     testServer.URL,
		Healthy: true,
	})

	// 执行测试用例
	tests := []struct {
		name           string
		path           string
		expectedStatus int
		timeout        time.Duration
	}{
		{
			name:           "Basic Proxy",
			path:           "/test/api",
			expectedStatus: http.StatusOK,
			timeout:        time.Second,
		},
		{
			name:           "Timeout Test",
			path:           "/test/slow",
			expectedStatus: http.StatusGatewayTimeout,
			timeout:        100 * time.Millisecond,
		},
		// 添加更多测试用例
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 执行测试
			resp, err := executeTest(tt.path, tt.timeout)
			if err != nil {
				t.Fatalf("Test failed: %v", err)
			}
			if resp.StatusCode != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, resp.StatusCode)
			}
		})
	}
}

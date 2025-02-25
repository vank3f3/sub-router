package handler

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"sub-router/internal/config"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestProxyHandler(t *testing.T) {
	// 设置测试配置
	config.GlobalConfig = config.Config{
		APIMappings: map[string]string{
			"test": "http://localhost:8888",
		},
	}

	// 创建测试服务器
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Test", "test")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test response"))
	}))
	defer backend.Close()

	// 更新测试配置
	config.GlobalConfig.APIMappings["test"] = backend.URL

	// 创建测试请求
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, r := gin.CreateTestContext(w)

	// 设置路由
	r.Any("/:service/*path", ProxyHandler)

	// 创建请求
	req := httptest.NewRequest("GET", "/test/api/data", nil)
	c.Request = req

	// 设置路由参数
	c.Params = gin.Params{
		{Key: "service", Value: "test"},
		{Key: "path", Value: "/api/data"},
	}

	// 执行处理器
	ProxyHandler(c)

	// 验证响应
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	if w.Header().Get("X-Test") != "test" {
		t.Error("Expected X-Test header to be set")
	}

	if body := w.Body.String(); body != "test response" {
		t.Errorf("Expected body %q, got %q", "test response", body)
	}
}

func TestProxyHandlerWithBody(t *testing.T) {
	// 设置测试配置
	config.GlobalConfig = config.Config{
		APIMappings: map[string]string{
			"test": "http://localhost:8888",
		},
	}

	// 创建测试服务器
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		w.Write(body) // 回显请求体
	}))
	defer backend.Close()

	// 更新测试配置
	config.GlobalConfig.APIMappings["test"] = backend.URL

	// 创建测试请求
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, r := gin.CreateTestContext(w)

	// 设置路由
	r.Any("/:service/*path", ProxyHandler)

	// 创建带请求体的请求
	requestBody := []byte(`{"test":"data"}`)
	req := httptest.NewRequest("POST", "/test/api/data",
		bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req

	// 设置路由参数
	c.Params = gin.Params{
		{Key: "service", Value: "test"},
		{Key: "path", Value: "/api/data"},
	}

	// 执行处理器
	ProxyHandler(c)

	// 验证响应
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	if body := w.Body.String(); body != string(requestBody) {
		t.Errorf("Expected body %q, got %q", string(requestBody), body)
	}
}

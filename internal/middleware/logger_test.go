package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func TestLogger(t *testing.T) {
	// 创建测试日志记录器
	logger, _ := zap.NewDevelopment()

	// 创建测试上下文
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// 创建请求
	req := httptest.NewRequest("GET", "/test", nil)
	c.Request = req

	// 执行中间件
	handler := Logger(logger)
	handler(c)

	// 验证响应
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}
}

func TestLoggerSkipPaths(t *testing.T) {
	// 创建测试日志记录器
	logger, _ := zap.NewDevelopment()

	// 测试跳过的路径
	skipPaths := []string{
		"/metrics",
		"/health",
		"/favicon.ico",
	}

	for _, path := range skipPaths {
		// 创建测试上下文
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// 创建请求
		req := httptest.NewRequest("GET", path, nil)
		c.Request = req

		// 执行中间件
		handler := Logger(logger)
		handler(c)

		// 验证响应
		if w.Code != http.StatusOK {
			t.Errorf("Expected status code %d for path %s, got %d",
				http.StatusOK, path, w.Code)
		}
	}
}

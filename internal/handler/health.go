package handler

import (
	"context"
	"net/http"
	"time"

	"sub-router/internal/config"

	"github.com/gin-gonic/gin"
)

// HealthStatus 健康状态
type HealthStatus struct {
	Status  string                 `json:"status"`
	Checks  map[string]CheckResult `json:"checks"`
	Version string                 `json:"version"`
}

// CheckResult 检查结果
type CheckResult struct {
	Status    string        `json:"status"`
	Message   string        `json:"message,omitempty"`
	Duration  time.Duration `json:"duration,omitempty"`
	Timestamp time.Time     `json:"timestamp"`
}

// HealthCheck 处理健康检查请求
func HealthCheck(c *gin.Context) {
	c.String(200, "Service is running!")
}

// RobotsHandler 处理 robots.txt 请求
func RobotsHandler(c *gin.Context) {
	c.Header("Content-Type", "text/plain")
	c.String(200, "User-agent: *\nDisallow: /")
}

// DetailedHealthCheck returns the health status of the service
func DetailedHealthCheck(c *gin.Context) {
	healthStatus := make(map[string]interface{})

	// Check API health
	healthStatus["api"] = "up" // Assuming API is always up for this example

	// Check proxy health only if enabled
	if config.GlobalConfig.Proxy.Enabled {
		// Perform proxy health check logic here
		// For example, you might want to ping the proxy server
		healthStatus["proxy"] = "up" // Replace with actual check
	} else {
		healthStatus["proxy"] = "not enabled"
	}

	c.JSON(http.StatusOK, healthStatus)
}

// performHealthCheck 执行健康检查
func performHealthCheck(name string, timeout time.Duration) CheckResult {
	start := time.Now()

	// 创建带超时的上下文
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// 根据不同的检查类型执行不同的检查
	switch name {
	case "proxy":
		return checkProxy(ctx)
	case "api":
		return checkAPI(ctx)
	default:
		return CheckResult{
			Status:    "unknown",
			Message:   "Unknown check type",
			Duration:  time.Since(start),
			Timestamp: time.Now(),
		}
	}
}

// checkProxy 检查代理服务是否正常
func checkProxy(ctx context.Context) CheckResult {
	start := time.Now()

	// 创建 HTTP 客户端
	client := &http.Client{}

	// 检查代理配置
	enabled, proxyURL := config.GetProxyConfig()
	if !enabled || proxyURL == "" {
		return CheckResult{
			Status:    "warning",
			Message:   "Proxy not configured",
			Duration:  time.Since(start),
			Timestamp: time.Now(),
		}
	}

	// 尝试连接代理
	req, err := http.NewRequestWithContext(ctx, "HEAD", proxyURL, nil)
	if err != nil {
		return CheckResult{
			Status:    "error",
			Message:   "Failed to create request: " + err.Error(),
			Duration:  time.Since(start),
			Timestamp: time.Now(),
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		return CheckResult{
			Status:    "error",
			Message:   "Failed to connect to proxy: " + err.Error(),
			Duration:  time.Since(start),
			Timestamp: time.Now(),
		}
	}
	defer resp.Body.Close()

	return CheckResult{
		Status:    "ok",
		Duration:  time.Since(start),
		Timestamp: time.Now(),
	}
}

// checkAPI 检查 API 服务是否正常
func checkAPI(ctx context.Context) CheckResult {
	start := time.Now()

	// 这里可以添加具体的 API 检查逻辑
	// 例如：检查一些关键的 API 端点是否可访问

	select {
	case <-ctx.Done():
		return CheckResult{
			Status:    "error",
			Message:   "API check timeout",
			Duration:  time.Since(start),
			Timestamp: time.Now(),
		}
	case <-time.After(100 * time.Millisecond): // 模拟检查耗时
		return CheckResult{
			Status:    "ok",
			Duration:  time.Since(start),
			Timestamp: time.Now(),
		}
	}
}

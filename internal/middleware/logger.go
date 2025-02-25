package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Logger 日志中间件
func Logger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 跳过特定路径的日志记录
		if shouldSkipLogging(c.Request.URL.Path) {
			c.Next()
			return
		}

		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		traceID := c.GetString("trace_id")

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()

		// 记录请求日志
		logger.Info("request",
			zap.Int("status", status),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.Duration("latency", latency),
			zap.String("ip", c.ClientIP()),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.String("trace_id", traceID),
		)
	}
}

// shouldSkipLogging 判断是否跳过日志记录
func shouldSkipLogging(path string) bool {
	// 跳过特定路径
	switch path {
	case "/metrics",
		"/health",
		"/favicon.ico":
		return true
	}
	return false
}

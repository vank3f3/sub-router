package middleware

import (
	"sub-router/internal/config"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Tracing 请求追踪中间件
func Tracing() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !config.GlobalConfig.Tracing.Enabled {
			c.Next()
			return
		}

		// 获取或生成追踪 ID
		traceID := c.GetHeader(config.GlobalConfig.Tracing.HeaderName)
		if traceID == "" {
			traceID = uuid.New().String()
		}

		// 设置追踪 ID
		c.Set("trace_id", traceID)
		c.Header(config.GlobalConfig.Tracing.HeaderName, traceID)

		c.Next()
	}
}

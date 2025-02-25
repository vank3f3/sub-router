package middleware

import (
	"sub-router/pkg/errors"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// RateLimit 限流中间件
func RateLimit(r rate.Limit, b int) gin.HandlerFunc {
	limiter := rate.NewLimiter(r, b)
	return func(c *gin.Context) {
		if !limiter.Allow() {
			c.AbortWithError(429,
				errors.New(errors.ErrorTypeRateLimit, "too many requests", 429))
			return
		}
		c.Next()
	}
}

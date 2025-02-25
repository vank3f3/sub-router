package middleware

import (
	"context"
	"net/http"
	"time"

	"sub-router/pkg/errors"

	"github.com/gin-gonic/gin"
)

// Timeout 超时中间件
func Timeout(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()

		c.Request = c.Request.WithContext(ctx)

		done := make(chan struct{})
		go func() {
			c.Next()
			done <- struct{}{}
		}()

		select {
		case <-done:
			return
		case <-ctx.Done():
			c.AbortWithError(http.StatusGatewayTimeout,
				errors.New(errors.ErrorTypeInternal, "request timeout", http.StatusGatewayTimeout))
		}
	}
}

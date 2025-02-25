package errors

import (
	"fmt"
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

// RecoveryHandler 统一的错误恢复处理
func RecoveryHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 记录堆栈信息
				stack := debug.Stack()

				// 构造错误响应
				errResp := &ErrorResponse{
					Type:    ErrorTypeInternal,
					Code:    500,
					Message: fmt.Sprintf("Internal Server Error: %v", err),
					Stack:   string(stack),
					TraceID: c.GetString("trace_id"),
				}

				// 返回错误响应
				c.JSON(500, errResp)

				// 终止后续处理
				c.Abort()
			}
		}()
		c.Next()
	}
}

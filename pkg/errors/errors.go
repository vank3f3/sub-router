package errors

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/pkg/errors"
)

// ErrorType 错误类型
type ErrorType string

const (
	ErrorTypeInternal   ErrorType = "INTERNAL_ERROR"
	ErrorTypeValidation ErrorType = "VALIDATION_ERROR"
	ErrorTypeProxy      ErrorType = "PROXY_ERROR"
	ErrorTypeConfig     ErrorType = "CONFIG_ERROR"
	ErrorTypePermission ErrorType = "PERMISSION_ERROR"
	ErrorTypeRateLimit  ErrorType = "RATE_LIMIT_ERROR"
	ErrorTypeAuth       ErrorType = "AUTH_ERROR"
	ErrorTypeNetwork    ErrorType = "NETWORK_ERROR"
	ErrorTypeTimeout    ErrorType = "TIMEOUT_ERROR"
	ErrorTypeThirdParty ErrorType = "THIRD_PARTY_ERROR"
)

// ErrorResponse 统一的错误响应结构
type ErrorResponse struct {
	Type    ErrorType `json:"type"`
	Code    int       `json:"code"`
	Message string    `json:"message"`
	TraceID string    `json:"trace_id,omitempty"`
	Stack   string    `json:"stack,omitempty"`
}

// APIError 自定义错误结构
type APIError struct {
	Type    ErrorType
	Message string
	Code    int
	err     error
	stack   string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

// New 创建新的 APIError
func New(errType ErrorType, message string, code int) *APIError {
	err := &APIError{
		Type:    errType,
		Message: message,
		Code:    code,
	}
	err.stack = getStack()
	return err
}

// Wrap 包装已有错误
func Wrap(err error, errType ErrorType, message string, code int) *APIError {
	return &APIError{
		Type:    errType,
		Message: message,
		Code:    code,
		err:     errors.Wrap(err, message),
		stack:   getStack(),
	}
}

// getStack 获取堆栈信息
func getStack() string {
	var stack strings.Builder
	for i := 2; i < 7; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		fn := runtime.FuncForPC(pc)
		stack.WriteString(fmt.Sprintf("%s:%d %s\n", file, line, fn.Name()))
	}
	return stack.String()
}

// ToResponse 转换为响应格式
func (e *APIError) ToResponse(traceID string) ErrorResponse {
	return ErrorResponse{
		Type:    e.Type,
		Code:    e.Code,
		Message: e.Message,
		TraceID: traceID,
		Stack:   e.stack,
	}
}

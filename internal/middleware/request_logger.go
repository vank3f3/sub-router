package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

// RequestLog 结构体用于记录请求日志
type RequestLog struct {
	Time          string `json:"time"`
	IP            string `json:"ip"`
	Authorization string `json:"authorization,omitempty"`
	Method        string `json:"method"`
	Path          string `json:"path"`
	Body          string `json:"body,omitempty"`
	TraceID       string `json:"trace_id,omitempty"`
}

// RequestLogger 中间件用于记录请求和响应信息
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {

		// 记录请求体
		bodyBytes, _ := io.ReadAll(c.Request.Body)
		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes)) // 重新设置请求体

		// 记录请求头
		logRequestHeaders(c, bodyBytes)

		// 处理请求
		c.Next()
	}
}

// logRequestHeaders 记录请求头
func logRequestHeaders(c *gin.Context, body []byte) {
	logFile, err := os.OpenFile("logs/requests.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Println("无法打开日志文件:", err)
		return
	}
	defer logFile.Close()

	// 压缩 JSON 请求体
	var compressedBody bytes.Buffer
	if json.Valid(body) {
		// 如果是有效的 JSON，进行压缩
		if err := json.Compact(&compressedBody, body); err == nil {
			body = compressedBody.Bytes()
		}
	}

	// 创建请求日志
	requestLog := RequestLog{
		Time: time.Now().Format(time.RFC3339), // 使用当前时间戳
		IP:   c.ClientIP(),
		//Authorization: c.Request.Header.Get("Authorization"),
		Method:  c.Request.Method,
		Path:    c.Request.URL.Path,
		TraceID: c.Request.Header.Get("X-Request-ID"), // 确保 trace_id 被记录
		Body:    string(body),
	}

	// 将请求日志转换为 JSON
	logData, _ := json.Marshal(requestLog)
	logFile.WriteString(string(logData) + "\n")
}

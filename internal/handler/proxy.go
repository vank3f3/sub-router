package handler

import (
	"io"
	"net/http"
	"net/url"
	"strings"
	"sub-router/internal/config"

	"github.com/gin-gonic/gin"
	"golang.org/x/net/proxy"
)

// ProxyHandler 处理代理请求
func ProxyHandler(c *gin.Context) {
	// 获取目标服务和路径
	service := c.Param("service")
	path := c.Param("path")

	// 获取目标基础URL
	baseURL, exists := config.GetAPIMapping(service)
	if !exists {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	// 构建目标URL
	targetURL := baseURL
	if path != "" {
		// 确保path不以/开头
		if strings.HasPrefix(path, "/") {
			path = path[1:]
		}
		// 确保baseURL以/结尾
		if !strings.HasSuffix(baseURL, "/") {
			targetURL += "/"
		}
		targetURL += path
	}

	// 添加查询参数
	if c.Request.URL.RawQuery != "" {
		targetURL += "?" + c.Request.URL.RawQuery
	}

	// 创建新的请求
	req, err := http.NewRequestWithContext(c.Request.Context(), c.Request.Method, targetURL, c.Request.Body)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// 复制请求头
	copyHeaders(c.Request.Header, req.Header)

	// 获取HTTP客户端
	client := getHTTPClient()

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		c.AbortWithStatus(http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	// 设置响应头
	copyHeaders(resp.Header, c.Writer.Header())

	// 设置状态码
	c.Status(resp.StatusCode)

	// 转发响应体
	io.Copy(c.Writer, resp.Body)
}

// copyHeaders 复制HTTP头
func copyHeaders(src, dst http.Header) {
	for key, values := range src {
		if isAllowedHeader(key) {
			for _, value := range values {
				dst.Add(key, value)
			}
		}
	}
}

// isAllowedHeader 检查是否是允许的请求头
func isAllowedHeader(header string) bool {
	header = strings.ToLower(header)
	switch header {
	case "content-length", "connection", "transfer-encoding":
		return false
	}
	return true
}

// getHTTPClient 获取HTTP客户端
func getHTTPClient() *http.Client {
	// 检查是否启用代理
	enabled, proxyURL := config.GetProxyConfig()
	if !enabled || proxyURL == "" {
		return http.DefaultClient
	}

	// 解析代理URL
	parsedURL, err := url.Parse(proxyURL)
	if err != nil {
		return http.DefaultClient
	}

	// 根据代理类型创建Transport
	transport := &http.Transport{}

	if parsedURL.Scheme == "socks5" {
		// SOCKS5代理
		dialer, err := proxy.SOCKS5("tcp", parsedURL.Host, nil, proxy.Direct)
		if err != nil {
			return http.DefaultClient
		}
		transport.Dial = dialer.Dial
	} else {
		// HTTP/HTTPS代理
		transport.Proxy = http.ProxyURL(parsedURL)
	}

	return &http.Client{Transport: transport}
}

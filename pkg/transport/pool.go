package transport

import (
	"crypto/tls"
	"net"
	"net/http"
	"sync"
	"time"

	"sub-router/internal/config"
)

// Pool HTTP连接池管理器
type Pool struct {
	config config.TransportPoolConfig
	client *http.Client
	mu     sync.RWMutex

	// 连接池优化
	MaxConnsPerHost       int
	MaxIdleConnsPerHost   int
	IdleConnTimeout       time.Duration
	TLSHandshakeTimeout   time.Duration
	ExpectContinueTimeout time.Duration
	ResponseHeaderTimeout time.Duration
}

// NewPool 创建新的连接池
func NewPool(config config.TransportPoolConfig) *Pool {
	p := &Pool{
		config:                config,
		MaxConnsPerHost:       100,
		MaxIdleConnsPerHost:   10,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		ResponseHeaderTimeout: 30 * time.Second,
	}

	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		ForceAttemptHTTP2:      true,
		MaxIdleConns:           config.MaxIdleConns,
		MaxIdleConnsPerHost:    config.MaxIdleConnsPerHost,
		IdleConnTimeout:        config.IdleConnTimeout,
		MaxConnsPerHost:        config.MaxIdleConnsPerHost * 2,
		ExpectContinueTimeout:  1 * time.Second,
		MaxResponseHeaderBytes: 4 * 1024, // 4KB
	}

	p.client = &http.Client{
		Transport: transport,
		Timeout:   30 * time.Second,
	}

	p.optimizeConnPool()
	return p
}

// Client 获取HTTP客户端
func (p *Pool) Client() *http.Client {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.client
}

// UpdateConfig 更新连接池配置
func (p *Pool) UpdateConfig(config config.TransportPoolConfig) {
	p.mu.Lock()
	defer p.mu.Unlock()

	transport := p.client.Transport.(*http.Transport)
	transport.MaxIdleConns = config.MaxIdleConns
	transport.MaxIdleConnsPerHost = config.MaxIdleConnsPerHost
	transport.IdleConnTimeout = config.IdleConnTimeout

	if config.InsecureSkipVerify {
		transport.TLSClientConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
	}

	p.config = config
}

// CloseIdleConnections 关闭所有空闲连接
func (p *Pool) CloseIdleConnections() {
	p.mu.RLock()
	defer p.mu.RUnlock()
	p.client.Transport.(*http.Transport).CloseIdleConnections()
}

// 优化连接池配置
func (p *Pool) optimizeConnPool() {
	transport := p.client.Transport.(*http.Transport)

	// 设置最大连接数
	transport.MaxConnsPerHost = p.MaxConnsPerHost
	transport.MaxIdleConnsPerHost = p.MaxIdleConnsPerHost

	// 设置超时
	transport.IdleConnTimeout = p.IdleConnTimeout
	transport.TLSHandshakeTimeout = p.TLSHandshakeTimeout
	transport.ExpectContinueTimeout = p.ExpectContinueTimeout
	transport.ResponseHeaderTimeout = p.ResponseHeaderTimeout

	// 启用 HTTP/2
	transport.ForceAttemptHTTP2 = true
}

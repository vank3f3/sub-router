package transport

import "sub-router/internal/config"

var (
	// GlobalPool 全局连接池实例
	GlobalPool *Pool
)

// InitGlobalPool 初始化全局连接池
func InitGlobalPool(config config.TransportPoolConfig) {
	GlobalPool = NewPool(config)
}

// CloseGlobalPool 关闭全局连接池
func CloseGlobalPool() {
	if GlobalPool != nil {
		GlobalPool.CloseIdleConnections()
	}
}

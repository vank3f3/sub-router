package websocket

import (
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// Config WebSocket配置
type Config struct {
	// 握手超时时间
	HandshakeTimeout time.Duration
	// 缓冲区大小
	BufferSize int
	// 允许的子协议
	Subprotocols []string
	// 允许的源
	AllowedOrigins []string
}

// Proxy WebSocket代理
type Proxy struct {
	config   Config
	upgrader websocket.Upgrader
}

// NewProxy 创建新的WebSocket代理
func NewProxy(config Config) *Proxy {
	return &Proxy{
		config: config,
		upgrader: websocket.Upgrader{
			HandshakeTimeout: config.HandshakeTimeout,
			ReadBufferSize:   config.BufferSize,
			WriteBufferSize:  config.BufferSize,
			Subprotocols:     config.Subprotocols,
			CheckOrigin: func(r *http.Request) bool {
				origin := r.Header.Get("Origin")
				if len(config.AllowedOrigins) == 0 {
					return true
				}
				for _, allowed := range config.AllowedOrigins {
					if allowed == "*" || allowed == origin {
						return true
					}
				}
				return false
			},
		},
	}
}

// ProxyHandler 处理WebSocket代理请求
func (p *Proxy) ProxyHandler(c *gin.Context, targetURL string) {
	// 升级HTTP连接为WebSocket
	conn, err := p.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	// 解析目标URL
	u, err := url.Parse(targetURL)
	if err != nil {
		return
	}

	// 修改协议
	if u.Scheme == "http" {
		u.Scheme = "ws"
	} else if u.Scheme == "https" {
		u.Scheme = "wss"
	}

	// 连接目标WebSocket服务器
	targetConn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return
	}
	defer targetConn.Close()

	// 双向转发数据
	errChan := make(chan error, 2)

	// 客户端 -> 服务器
	go p.transfer(conn, targetConn, errChan)
	// 服务器 -> 客户端
	go p.transfer(targetConn, conn, errChan)

	// 等待任一方向出错
	<-errChan
}

// transfer 在两个WebSocket连接之间传输数据
func (p *Proxy) transfer(src, dst *websocket.Conn, errChan chan error) {
	for {
		messageType, message, err := src.ReadMessage()
		if err != nil {
			errChan <- err
			return
		}

		err = dst.WriteMessage(messageType, message)
		if err != nil {
			errChan <- err
			return
		}
	}
}

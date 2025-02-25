package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Options 服务器选项
type Options struct {
	Addr           string
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	MaxHeaderBytes int
	Logger         *zap.Logger
}

// Server HTTP服务器
type Server struct {
	*http.Server
	logger *zap.Logger
	router *gin.Engine
}

// NewServer 创建新的服务器
func NewServer(router *gin.Engine, opts Options) *Server {
	return &Server{
		Server: &http.Server{
			Addr:           opts.Addr,
			Handler:        router,
			ReadTimeout:    opts.ReadTimeout,
			WriteTimeout:   opts.WriteTimeout,
			MaxHeaderBytes: opts.MaxHeaderBytes,
		},
		logger: opts.Logger,
		router: router,
	}
}

// Start 启动服务器
func (s *Server) Start() error {
	// 创建错误通道
	errChan := make(chan error, 1)

	// 启动服务器
	go func() {
		s.logger.Info("Starting server", zap.String("addr", s.Addr))
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- err
		}
	}()

	// 监听信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errChan:
		return err
	case <-quit:
		s.logger.Info("Shutting down server...")
		return s.Shutdown()
	}
}

// Shutdown 优雅关闭服务器
func (s *Server) Shutdown() error {
	// 创建关闭上下文
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 关闭服务器
	if err := s.Server.Shutdown(ctx); err != nil {
		s.logger.Error("Server forced to shutdown", zap.Error(err))
		return err
	}

	s.logger.Info("Server exiting")
	return nil
}

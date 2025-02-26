package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"golang.org/x/time/rate"

	"sub-router/internal/config"
	"sub-router/internal/handler"
	"sub-router/internal/middleware"
	"sub-router/pkg/transport"
)

func main() {
	// 初始化日志
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	// 创建日志目录
	if err := os.MkdirAll("logs", os.ModePerm); err != nil {
		log.Fatalf("无法创建日志目录: %v", err)
	}

	// 加载配置
	if err := config.LoadConfig(); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 获取服务器配置
	port, timeout, rateLimit := config.GetServerConfig()

	// 初始化全局连接池
	transport.InitGlobalPool(config.GetTransportConfig())
	defer transport.CloseGlobalPool()

	// 创建 gin 引擎
	ginMode := config.GlobalConfig.Server.GinMode // 读取 GIN_MODE
	gin.SetMode(ginMode)                          // 设置 GIN_MODE
	r := gin.New()

	// 添加中间件（注意顺序）
	r.Use(middleware.RequestLogger()) // 添加请求日志中间件
	if ginMode != "release" {
		r.Use(gin.Logger()) // 仅在非 release 模式下使用日志中间件
	}
	r.Use(gin.Recovery())              // 错误恢复
	r.Use(middleware.Tracing())        // 请求追踪
	r.Use(middleware.Logger(logger))   // 自定义日志记录
	r.Use(middleware.Security())       // 安全头
	r.Use(middleware.IPControl())      // IP 控制
	r.Use(middleware.Metrics())        // 指标收集
	r.Use(middleware.Timeout(timeout)) // 超时控制
	r.Use(middleware.RateLimit(rate.Limit(rateLimit.RequestsPerSecond), rateLimit.Burst))

	// 监控路由
	if config.GlobalConfig.Monitoring.Metrics.Enabled {
		r.GET(config.GlobalConfig.Monitoring.Metrics.Path, gin.WrapH(promhttp.Handler()))
	}

	// 健康检查路由
	if config.GlobalConfig.Monitoring.Health.Enabled {
		r.GET(config.GlobalConfig.Monitoring.Health.DetailedPath, handler.DetailedHealthCheck)
	}

	// 基础路由
	r.GET("/", handler.HealthCheck)
	r.GET("/index.html", handler.HealthCheck)
	r.GET("/robots.txt", handler.RobotsHandler)

	// API 代理路由（确保路径处理正确）
	r.Any("/:service/*path", handler.ProxyHandler)

	// 启动服务器
	log.Printf("Server starting on port %d...", port)
	if err := r.Run(fmt.Sprintf(":%d", port)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

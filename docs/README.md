# Sub-Router

高性能的 API 代理和负载均衡服务器

## 功能特性

- 动态路由和代理
- 负载均衡
- 熔断保护
- 监控指标
- 性能优化

## 快速开始

### 安装

```bash
# 克隆仓库
git clone https://github.com/yourusername/sub-router.git
cd sub-router

# 安装依赖
go mod download

# 构建
go build -o sub-router

# 运行
./sub-router
```

### Docker 部署

```bash
# 构建镜像
docker build -t sub-router .

# 运行容器
docker run -d -p 8080:8080 -v ./configs:/app/configs sub-router
```

## 配置说明

配置文件位于 `configs/config.yaml`：

### 服务器配置
```yaml
server:
  port: 8080
  timeout: 30s
  rate_limit:
    requests_per_second: 100
    burst: 200
  gin_mode: "release" # Set to "release" for production
```

### API 映射配置
```yaml
api_mappings:
  discord: "https://discord.com/api"
  telegram: "https://api.telegram.org"
  openai: "https://api.openai.com"
```

### 代理配置
```yaml
proxy:
  enabled: true
  url: "socks5://127.0.0.1:7890"
```

### 监控配置
```yaml
monitoring:
  metrics:
    enabled: true
    path: "/metrics"
  health:
    enabled: true
    detailed_path: "/health"
```

## API 文档

### 代理请求
- 路径: `/:service/*path`
- 方法: 支持所有 HTTP 方法
- 示例: 
  ```
  GET /openai/v1/chat/completions
  POST /discord/api/webhooks
  ```

### 健康检查
- 基础检查: `GET /health`
- 详细检查: `GET /health/detail`

### 监控指标
- Prometheus 指标: `GET /metrics`

## 性能指标

### 基准测试结果
- 并发请求: 5000 QPS
- 平均延迟: < 100ms
- 内存占用: < 100MB

### 限制
- 最大连接数: 10000
- 单个请求超时: 30s
- 最大请求体积: 10MB

## 部署指南

### 系统要求
- Go 1.21+
- 2GB+ RAM
- 现代 Linux/Unix 系统

### 生产环境配置
1. 调整系统限制
```bash
# /etc/security/limits.conf
* soft nofile 65535
* hard nofile 65535
```

2. 网络优化
```bash
# /etc/sysctl.conf
net.ipv4.tcp_fin_timeout = 30
net.ipv4.tcp_keepalive_time = 1200
net.core.somaxconn = 65535
```

### 监控集成
- Prometheus + Grafana 面板
- ELK 日志收集
- 链路追踪

### 高可用部署
- 多实例负载均衡
- 健康检查和自动恢复
- 容器化部署

## 常见问题

### 性能调优
1. 连接池配置
2. 内存管理
3. 超时控制

### 故障排除
1. 日志分析
2. 指标监控
3. 链路追踪

## 贡献指南

1. Fork 项目
2. 创建特性分支
3. 提交变更
4. 发起 Pull Request

## 许可证

MIT License
package config

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

// ServerConfig 服务器配置
type ServerConfig struct {
	Port      int           `mapstructure:"port"`
	Timeout   time.Duration `mapstructure:"timeout"`
	RateLimit struct {
		RequestsPerSecond float64 `mapstructure:"requests_per_second"`
		Burst             int     `mapstructure:"burst"`
	} `mapstructure:"rate_limit"`
}

// ProxyConfig 代理配置
type ProxyConfig struct {
	Enabled bool   `mapstructure:"enabled"`
	URL     string `mapstructure:"url"`
}

// SecurityConfig 安全配置
type SecurityConfig struct {
	IPControl IPControlConfig `mapstructure:"ip_control"`
	BasicAuth BasicAuthConfig `mapstructure:"basic_auth"`
}

// IPControlConfig IP控制配置
type IPControlConfig struct {
	Enabled   bool     `mapstructure:"enabled"`
	Whitelist []string `mapstructure:"whitelist"`
	Blacklist []string `mapstructure:"blacklist"`
}

// BasicAuthConfig 基本认证配置
type BasicAuthConfig struct {
	Enabled     bool                   `mapstructure:"enabled"`
	Credentials []BasicAuthCredentials `mapstructure:"credentials"`
}

// BasicAuthCredentials 认证凭证
type BasicAuthCredentials struct {
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

// MonitoringConfig 监控配置
type MonitoringConfig struct {
	Metrics MetricsConfig `mapstructure:"metrics"`
	Health  HealthConfig  `mapstructure:"health"`
}

// MetricsConfig 指标配置
type MetricsConfig struct {
	Enabled bool   `mapstructure:"enabled"`
	Path    string `mapstructure:"path"`
}

// HealthConfig 健康检查配置
type HealthConfig struct {
	Enabled      bool          `mapstructure:"enabled"`
	DetailedPath string        `mapstructure:"detailed_path"`
	Checks       []HealthCheck `mapstructure:"checks"`
}

// HealthCheck 健康检查项
type HealthCheck struct {
	Name    string        `mapstructure:"name"`
	Timeout time.Duration `mapstructure:"timeout"`
}

// TracingConfig 追踪配置
type TracingConfig struct {
	Enabled    bool   `mapstructure:"enabled"`
	HeaderName string `mapstructure:"header_name"`
}

// TransportConfig 传输配置
type TransportConfig struct {
	MaxIdleConns        int           `mapstructure:"max_idle_conns"`
	MaxIdleConnsPerHost int           `mapstructure:"max_idle_conns_per_host"`
	IdleConnTimeout     time.Duration `mapstructure:"idle_conn_timeout"`
	MaxConnLifetime     time.Duration `mapstructure:"max_conn_lifetime"`
	TLSSkipVerify       bool          `mapstructure:"tls_skip_verify"`
}

// CompressionConfig 压缩配置
type CompressionConfig struct {
	Enabled bool   `mapstructure:"enabled"`
	Level   string `mapstructure:"level"`
}

// Config 总配置结构
type Config struct {
	APIMappings map[string]string `mapstructure:"api_mappings"`
	Server      ServerConfig      `mapstructure:"server"`
	Proxy       ProxyConfig       `mapstructure:"proxy"`
	Security    SecurityConfig    `mapstructure:"security"`
	Monitoring  MonitoringConfig  `mapstructure:"monitoring"`
	Tracing     TracingConfig     `mapstructure:"tracing"`
	Transport   TransportConfig   `mapstructure:"transport"`
	Compression CompressionConfig `mapstructure:"compression"`
}

var GlobalConfig Config

// LoadConfig 加载配置文件
func LoadConfig() error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")         // 首先在根目录查找
	viper.AddConfigPath("./configs") // 然后在configs目录查找

	// 设置默认值
	setDefaultConfig()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("No config file found, using defaults")
		} else {
			return err
		}
	}

	// 解析配置到结构体
	if err := viper.Unmarshal(&GlobalConfig); err != nil {
		return err
	}

	// 监听配置文件变化
	viper.WatchConfig()

	return nil
}

// setDefaultConfig 设置默认配置
func setDefaultConfig() {
	// 服务器默认配置
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.timeout", "30s")
	viper.SetDefault("server.rate_limit.requests_per_second", 100)
	viper.SetDefault("server.rate_limit.burst", 200)

	// 传输层默认配置
	viper.SetDefault("transport.max_idle_conns", 100)
	viper.SetDefault("transport.max_idle_conns_per_host", 10)
	viper.SetDefault("transport.idle_conn_timeout", "90s")
	viper.SetDefault("transport.max_conn_lifetime", "4m")
	viper.SetDefault("transport.tls_skip_verify", false)

	// 压缩默认配置
	viper.SetDefault("compression.enabled", true)
	viper.SetDefault("compression.level", "default")
}

// GetServerConfig 获取服务器配置
func GetServerConfig() (port int, timeout time.Duration, rateLimit struct {
	RequestsPerSecond float64
	Burst             int
}) {
	return GlobalConfig.Server.Port,
		GlobalConfig.Server.Timeout,
		struct {
			RequestsPerSecond float64
			Burst             int
		}{
			RequestsPerSecond: GlobalConfig.Server.RateLimit.RequestsPerSecond,
			Burst:             GlobalConfig.Server.RateLimit.Burst,
		}
}

// GetProxyConfig 获取代理配置
func GetProxyConfig() (enabled bool, proxyURL string) {
	return GlobalConfig.Proxy.Enabled, GlobalConfig.Proxy.URL
}

// GetAPIMapping 获取 API 映射
func GetAPIMapping(service string) (string, bool) {
	baseURL, exists := GlobalConfig.APIMappings[service]
	return baseURL, exists
}

// TransportPoolConfig 连接池配置
type TransportPoolConfig struct {
	MaxIdleConns        int
	MaxIdleConnsPerHost int
	IdleConnTimeout     time.Duration
	MaxConnLifetime     time.Duration
	InsecureSkipVerify  bool
}

// GetTransportConfig 获取传输配置
func GetTransportConfig() TransportPoolConfig {
	return TransportPoolConfig{
		MaxIdleConns:        GlobalConfig.Transport.MaxIdleConns,
		MaxIdleConnsPerHost: GlobalConfig.Transport.MaxIdleConnsPerHost,
		IdleConnTimeout:     GlobalConfig.Transport.IdleConnTimeout,
		MaxConnLifetime:     GlobalConfig.Transport.MaxConnLifetime,
		InsecureSkipVerify:  GlobalConfig.Transport.TLSSkipVerify,
	}
}

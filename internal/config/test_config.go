package config

import "time"

// LoadTestConfig 加载测试配置
func LoadTestConfig() {
	GlobalConfig = Config{
		Server: ServerConfig{
			Port:    8080,
			Timeout: 30 * time.Second,
			RateLimit: struct {
				RequestsPerSecond float64 `mapstructure:"requests_per_second"`
				Burst             int     `mapstructure:"burst"`
			}{
				RequestsPerSecond: 100,
				Burst:             200,
			},
		},
		Monitoring: MonitoringConfig{
			Metrics: MetricsConfig{
				Enabled: true,
				Path:    "/metrics",
			},
			Health: HealthConfig{
				Enabled:      true,
				DetailedPath: "/health",
				Checks: []HealthCheck{
					{
						Name:    "proxy",
						Timeout: 5 * time.Second,
					},
				},
			},
		},
	}
}

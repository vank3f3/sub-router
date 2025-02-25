package config

import (
	"fmt"
)

// ValidateConfig 验证配置的合法性
func ValidateConfig(cfg *Config) error {
	// 验证服务器配置
	if err := validateServerConfig(cfg.Server); err != nil {
		return fmt.Errorf("server config: %w", err)
	}

	// 验证代理配置
	if err := validateProxyConfig(cfg.Proxy); err != nil {
		return fmt.Errorf("proxy config: %w", err)
	}

	// 验证监控配置
	if err := validateMonitoringConfig(cfg.Monitoring); err != nil {
		return fmt.Errorf("monitoring config: %w", err)
	}

	return nil
}

// validateServerConfig 验证服务器配置
func validateServerConfig(cfg ServerConfig) error {
	if cfg.Port <= 0 || cfg.Port > 65535 {
		return fmt.Errorf("invalid port: %d", cfg.Port)
	}
	if cfg.Timeout < 0 {
		return fmt.Errorf("invalid timeout: %v", cfg.Timeout)
	}
	return nil
}

// validateProxyConfig 验证代理配置
func validateProxyConfig(cfg ProxyConfig) error {
	if cfg.Enabled {
		if cfg.URL == "" {
			return fmt.Errorf("proxy enabled but URL is empty")
		}
	}
	return nil
}

// validateMonitoringConfig 验证监控配置
func validateMonitoringConfig(cfg MonitoringConfig) error {
	if cfg.Metrics.Enabled {
		if cfg.Metrics.Path == "" {
			return fmt.Errorf("metrics enabled but path is empty")
		}
	}
	if cfg.Health.Enabled {
		if cfg.Health.DetailedPath == "" {
			return fmt.Errorf("health check enabled but path is empty")
		}
		for _, check := range cfg.Health.Checks {
			if check.Timeout <= 0 {
				return fmt.Errorf("invalid health check timeout: %v", check.Timeout)
			}
		}
	}
	return nil
}

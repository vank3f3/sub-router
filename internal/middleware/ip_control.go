package middleware

import (
	"net"
	"strings"

	"sub-router/internal/config"
	"sub-router/pkg/errors"

	"github.com/gin-gonic/gin"
)

// IPControl IP 控制中间件
func IPControl() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !config.GlobalConfig.Security.IPControl.Enabled {
			c.Next()
			return
		}

		clientIP := c.ClientIP()
		ip := net.ParseIP(clientIP)
		if ip == nil {
			c.AbortWithStatusJSON(403,
				errors.New(errors.ErrorTypePermission, "Invalid IP address", 403).
					ToResponse(c.GetString("trace_id")))
			return
		}

		// 检查黑名单
		for _, blackIP := range config.GlobalConfig.Security.IPControl.Blacklist {
			if isIPInRange(ip, blackIP) {
				c.AbortWithStatusJSON(403,
					errors.New(errors.ErrorTypePermission, "IP blocked", 403).
						ToResponse(c.GetString("trace_id")))
				return
			}
		}

		// 检查白名单
		if len(config.GlobalConfig.Security.IPControl.Whitelist) > 0 {
			allowed := false
			for _, whiteIP := range config.GlobalConfig.Security.IPControl.Whitelist {
				if isIPInRange(ip, whiteIP) {
					allowed = true
					break
				}
			}
			if !allowed {
				c.AbortWithStatusJSON(403,
					errors.New(errors.ErrorTypePermission, "IP not allowed", 403).
						ToResponse(c.GetString("trace_id")))
				return
			}
		}

		c.Next()
	}
}

// isIPInRange 检查 IP 是否在指定范围内
func isIPInRange(ip net.IP, ipRange string) bool {
	if strings.Contains(ipRange, "/") {
		_, ipNet, err := net.ParseCIDR(ipRange)
		if err != nil {
			return false
		}
		return ipNet.Contains(ip)
	}
	return ip.String() == ipRange
}

package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	// 请求延迟分布
	RequestLatency = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "request_latency_seconds",
			Help:    "Request latency distribution",
			Buckets: []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5},
		},
		[]string{"service", "method", "status"},
	)

	// 后端健康状态
	BackendHealth = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "backend_health_status",
			Help: "Backend health status (1 for healthy, 0 for unhealthy)",
		},
		[]string{"backend"},
	)

	// 熔断器状态
	CircuitBreakerStatus = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "circuit_breaker_status",
			Help: "Circuit breaker status (0: Closed, 1: Open, 2: Half-Open)",
		},
		[]string{"service"},
	)
)

func init() {
	prometheus.MustRegister(RequestLatency)
	prometheus.MustRegister(BackendHealth)
	prometheus.MustRegister(CircuitBreakerStatus)
}

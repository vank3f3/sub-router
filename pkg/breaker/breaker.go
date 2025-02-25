package breaker

import (
	"sync"
	"time"
)

// State 熔断器状态
type State int

const (
	StateClosed   State = iota // 关闭状态（正常）
	StateOpen                  // 开启状态（熔断）
	StateHalfOpen              // 半开状态（尝试恢复）
)

// Config 熔断器配置
type Config struct {
	ErrorThreshold   int           // 错误阈值
	SuccessThreshold int           // 成功阈值
	Timeout          time.Duration // 熔断超时时间
	MaxRequests      int           // 半开状态最大请求数
}

// CircuitBreaker 熔断器
type CircuitBreaker struct {
	state           State
	config          Config
	failures        int
	successes       int
	lastStateChange time.Time
	mu              sync.RWMutex
}

// NewCircuitBreaker 创建新的熔断器
func NewCircuitBreaker(config Config) *CircuitBreaker {
	return &CircuitBreaker{
		state:  StateClosed,
		config: config,
	}
}

// Allow 判断是否允许请求
func (cb *CircuitBreaker) Allow() bool {
	cb.mu.RLock()
	defer cb.mu.RUnlock()

	switch cb.state {
	case StateClosed:
		return true
	case StateOpen:
		if time.Since(cb.lastStateChange) > cb.config.Timeout {
			cb.toHalfOpen()
			return true
		}
		return false
	case StateHalfOpen:
		return cb.successes < cb.config.MaxRequests
	default:
		return false
	}
}

// Success 记录成功请求
func (cb *CircuitBreaker) Success() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	switch cb.state {
	case StateHalfOpen:
		cb.successes++
		if cb.successes >= cb.config.SuccessThreshold {
			cb.toClosed()
		}
	case StateClosed:
		cb.failures = 0
		cb.successes = 0
	}
}

// Failure 记录失败请求
func (cb *CircuitBreaker) Failure() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	switch cb.state {
	case StateClosed:
		cb.failures++
		if cb.failures >= cb.config.ErrorThreshold {
			cb.toOpen()
		}
	case StateHalfOpen:
		cb.toOpen()
	}
}

// 状态转换方法
func (cb *CircuitBreaker) toOpen() {
	cb.state = StateOpen
	cb.lastStateChange = time.Now()
	cb.failures = 0
	cb.successes = 0
}

func (cb *CircuitBreaker) toHalfOpen() {
	cb.state = StateHalfOpen
	cb.lastStateChange = time.Now()
	cb.failures = 0
	cb.successes = 0
}

func (cb *CircuitBreaker) toClosed() {
	cb.state = StateClosed
	cb.lastStateChange = time.Now()
	cb.failures = 0
	cb.successes = 0
}

// State 获取当前状态
func (cb *CircuitBreaker) State() State {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	return cb.state
}

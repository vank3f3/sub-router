package breaker

import (
	"testing"
	"time"
)

func TestCircuitBreaker(t *testing.T) {
	// 创建熔断器配置
	config := Config{
		ErrorThreshold:   3,
		SuccessThreshold: 2,
		Timeout:          1 * time.Second,
		MaxRequests:      2,
	}

	// 创建熔断器
	breaker := NewCircuitBreaker(config)

	// 测试初始状态
	if state := breaker.State(); state != StateClosed {
		t.Errorf("Expected initial state to be Closed, got %v", state)
	}

	// 测试错误阈值触发熔断
	for i := 0; i < config.ErrorThreshold; i++ {
		breaker.Failure()
	}
	if state := breaker.State(); state != StateOpen {
		t.Errorf("Expected state to be Open after failures, got %v", state)
	}

	// 测试熔断状态下拒绝请求
	if breaker.Allow() {
		t.Error("Expected request to be rejected in Open state")
	}

	// 等待超时时间
	time.Sleep(config.Timeout)

	// 测试半开状态
	if !breaker.Allow() {
		t.Error("Expected request to be allowed in HalfOpen state")
	}

	// 测试成功请求导致关闭
	for i := 0; i < config.SuccessThreshold; i++ {
		breaker.Success()
	}
	if state := breaker.State(); state != StateClosed {
		t.Errorf("Expected state to be Closed after successes, got %v", state)
	}
}

func TestCircuitBreakerEdgeCases(t *testing.T) {
	config := Config{
		ErrorThreshold:   2,
		SuccessThreshold: 1,
		Timeout:          100 * time.Millisecond,
		MaxRequests:      1,
	}
	breaker := NewCircuitBreaker(config)

	// 测试单个失败不触发熔断
	breaker.Failure()
	if state := breaker.State(); state != StateClosed {
		t.Error("Single failure should not trigger circuit breaker")
	}

	// 测试半开状态下的失败立即触发熔断
	breaker.Failure() // 触发熔断
	time.Sleep(config.Timeout)
	breaker.Allow() // 进入半开状态
	breaker.Failure()
	if state := breaker.State(); state != StateOpen {
		t.Error("Failure in half-open state should trigger immediate open")
	}

	// 测试并发安全性
	go func() {
		for i := 0; i < 100; i++ {
			breaker.Success()
		}
	}()
	go func() {
		for i := 0; i < 100; i++ {
			breaker.Failure()
		}
	}()
	time.Sleep(100 * time.Millisecond)
}

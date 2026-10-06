// Package health 提供应用健康状态检查的基础能力。
package health

import (
	"sync"
	"time"
)

// HealthChecker 是健康状态检查器，记录应用当前是否正常运行。
type HealthChecker struct {
	mu      sync.RWMutex
	healthy bool
	details map[string]string
}

// NewHealthChecker 创建一个健康检查器，初始状态为健康
func NewHealthChecker() *HealthChecker {
	return &HealthChecker{
		healthy: true,
		details: make(map[string]string),
	}
}

// IsHealthy 返回当前应用是否健康
func (c *HealthChecker) IsHealthy() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.healthy
}

// SetHealthy 设置健康状态
func (c *HealthChecker) SetHealthy(v bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.healthy = v
	if !v {
		c.details["since"] = time.Now().Format(time.RFC3339)
	}
}

// SetDetail 设置或更新一个健康状态细节
func (c *HealthChecker) SetDetail(key, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.details[key] = value
}

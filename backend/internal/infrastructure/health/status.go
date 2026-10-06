package health

import (
	"fmt"
	"maps"
	"time"
)

// HealthStatus 返回当前健康状态的快照。
type HealthStatus struct {
	Status    string            `json:"status"`
	Timestamp string            `json:"timestamp"`
	Details   map[string]string `json:"details,omitempty"`
	Uptime    string            `json:"uptime,omitempty"`
}

// Snapshot 返回当前健康状态的只读快照。
func (c *HealthChecker) Snapshot(startTime time.Time) HealthStatus {
	c.mu.RLock()
	defer c.mu.RUnlock()

	status := "ok"
	if !c.healthy {
		status = "unhealthy"
	}
	details := make(map[string]string, len(c.details))
	maps.Copy(details, c.details)

	uptime := time.Since(startTime).Truncate(time.Second).String()

	return HealthStatus{
		Status:    status,
		Timestamp: time.Now().Format(time.RFC3339),
		Details:   details,
		Uptime:    fmt.Sprintf("%s", uptime),
	}
}

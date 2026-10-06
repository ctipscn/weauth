package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/ctipscn/weauth/backend/internal/interface/http/response"

	"github.com/gin-gonic/gin"
)

// slidingWindow 滑动窗口计数器
type slidingWindow struct {
	start time.Time
	count int
}

// sitekeyBody 用于从请求体中提取 sitekey
type sitekeyBody struct {
	Sitekey string `json:"sitekey"`
}

// LimiterConfig 限流器配置
type LimiterConfig struct {
	TimeWindowSeconds int
	RequestLimit      int
}

// TrafficLimiter 基于内存的滑动窗口限流中间件。
// 根据请求体中的 sitekey 对每个站点进行独立限流。
type TrafficLimiter struct {
	mu      sync.Mutex
	windows map[string]*slidingWindow
	cfg     *LimiterConfig
}

// NewTrafficLimiter 创建 TrafficLimiter 实例
func NewTrafficLimiter(cfg *LimiterConfig) *TrafficLimiter {
	return &TrafficLimiter{
		windows: make(map[string]*slidingWindow),
		cfg:     cfg,
	}
}

// Handler 返回 gin 中间件处理器。
// 从 POST 请求体中提取 sitekey，对每个 sitekey 独立统计请求次数。
func (tl *TrafficLimiter) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 读取原始请求体
		body, err := c.GetRawData()
		if err != nil {
			c.JSON(http.StatusTooManyRequests, response.ErrorBaseResponse(
				http.StatusTooManyRequests, 0, "rate limit exceeded"))
			c.Abort()
			return
		}
		// 恢复请求体，不影响下游 handler
		c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

		// 尝试解析 sitekey
		var parsed sitekeyBody
		if err := json.Unmarshal(body, &parsed); err != nil || parsed.Sitekey == "" {
			// 无法识别站点标识则放行（避免影响非 PoW 端点）
			c.Next()
			return
		}

		windowDuration := time.Duration(tl.cfg.TimeWindowSeconds) * time.Second
		now := time.Now()

		tl.mu.Lock()
		w, exists := tl.windows[parsed.Sitekey]

		// 窗口过期或不存，重置
		if !exists || now.Sub(w.start) >= windowDuration {
			tl.windows[parsed.Sitekey] = &slidingWindow{
				start: now,
				count: 1,
			}
			tl.mu.Unlock()
			c.Next()
			return
		}

		// 窗口内计数递增
		w.count++
		if w.count > tl.cfg.RequestLimit {
			tl.mu.Unlock()
			c.JSON(http.StatusTooManyRequests, response.ErrorBaseResponse(
				http.StatusTooManyRequests, 0, "rate limit exceeded"))
			c.Abort()
			return
		}
		tl.mu.Unlock()
		c.Next()
	}
}

// Cleanup 定期清理过期窗口，防止内存泄漏。
// 建议在服务启动时以 goroutine 方式调用。
func (tl *TrafficLimiter) Cleanup(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for range ticker.C {
		windowDuration := time.Duration(tl.cfg.TimeWindowSeconds) * time.Second
		now := time.Now()

		tl.mu.Lock()
		for key, w := range tl.windows {
			if now.Sub(w.start) >= windowDuration {
				delete(tl.windows, key)
			}
		}
		tl.mu.Unlock()
	}
}

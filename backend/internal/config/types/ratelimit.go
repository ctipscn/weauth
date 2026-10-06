package types

import (
	"os"
	"sync"

	"github.com/ctipscn/weauth/backend/internal/config/utils"
)

// RateLimitConfig 限流器配置
type RateLimitConfig struct {
	TimeWindowSeconds int // 验证窗口（秒）
	RequestLimit      int // 窗口内最大请求次数
}

var (
	rlInstance *RateLimitConfig
	rlOnce     sync.Once
)

// RateLimit 返回全局单例 RateLimitConfig。
// 必须先调用 NewDefaultRateLimitConfig 初始化，否则返回 nil。
func RateLimit() *RateLimitConfig {
	return rlInstance
}

// ApplyFormFile 从 INI 文件的 [ratelimit] 节加载配置并合并到当前实例。
// 支持链式调用。
func (c *RateLimitConfig) ApplyFormFile(path string) (*RateLimitConfig, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return c, nil
	}
	sections, err := utils.ParseINI(path)
	if err != nil {
		return nil, err
	}
	if kv, ok := sections["ratelimit"]; ok {
		utils.ApplyInt(kv, "time_window_seconds", &c.TimeWindowSeconds)
		utils.ApplyInt(kv, "request_limit", &c.RequestLimit)
	}
	return c, nil
}

// ApplyFromENV 从 WEAUTH_ 前缀的环境变量加载配置并合并到当前实例。
// 支持链式调用。
func (c *RateLimitConfig) ApplyFromENV() *RateLimitConfig {
	utils.EnvInt("WEAUTH_RATELIMIT_TIME_WINDOW_SECONDS", &c.TimeWindowSeconds)
	utils.EnvInt("WEAUTH_RATELIMIT_REQUEST_LIMIT", &c.RequestLimit)
	return c
}

// NewDefaultRateLimitConfig 初始化并返回全局单例 RateLimitConfig。
// 后续调用直接返回同一实例。
func NewDefaultRateLimitConfig() *RateLimitConfig {
	rlOnce.Do(func() {
		rlInstance = &RateLimitConfig{
			TimeWindowSeconds: 30,
			RequestLimit:      500,
		}
	})
	return rlInstance
}

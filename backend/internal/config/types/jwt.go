package types

import (
	"os"
	"sync"

	"github.com/ctipscn/weauth/backend/internal/config/utils"
)

// JWTConfig JWT 配置
type JWTConfig struct {
	Secret      string
	ExpireHours int
}

var (
	jwtInstance *JWTConfig
	jwtOnce     sync.Once
)

// JWT 返回全局单例 JWTConfig。
// 必须先调用 NewDefaultJWTConfig 初始化，否则返回 nil。
func JWT() *JWTConfig {
	return jwtInstance
}

// ApplyFormFile 从 INI 文件的 [jwt] 节加载配置并合并到当前实例。
// 支持链式调用。
func (c *JWTConfig) ApplyFormFile(path string) (*JWTConfig, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return c, nil
	}
	sections, err := utils.ParseINI(path)
	if err != nil {
		return nil, err
	}
	if kv, ok := sections["jwt"]; ok {
		utils.ApplyString(kv, "secret", &c.Secret)
		utils.ApplyInt(kv, "expire_hours", &c.ExpireHours)
	}
	return c, nil
}

// ApplyFromENV 从 WEAUTH_ 前缀的环境变量加载配置并合并到当前实例。
// 支持链式调用。
func (c *JWTConfig) ApplyFromENV() *JWTConfig {
	utils.EnvString("WEAUTH_JWT_SECRET", &c.Secret)
	utils.EnvInt("WEAUTH_JWT_EXPIRE_HOURS", &c.ExpireHours)
	return c
}

// NewDefaultJWTConfig 初始化并返回全局单例 JWTConfig。
// 后续调用直接返回同一实例。
func NewDefaultJWTConfig() *JWTConfig {
	jwtOnce.Do(func() {
		jwtInstance = &JWTConfig{
			Secret:      "default-secret-please-change",
			ExpireHours: 24,
		}
	})
	return jwtInstance
}
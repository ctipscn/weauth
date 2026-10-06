package types

import (
	"os"
	"sync"

	"github.com/ctipscn/weauth/backend/internal/config/utils"
)

// CORSConfig CORS 跨域配置
type CORSConfig struct {
	AllowOrigins     []string
	AllowMethods     []string
	AllowHeaders     []string
	ExposeHeaders    []string
	AllowCredentials bool
}

var (
	corsInstance *CORSConfig
	corsOnce     sync.Once
)

// CORS 返回全局单例 CORSConfig。
// 必须先调用 NewDefaultCORSConfig 初始化，否则返回 nil。
func CORS() *CORSConfig {
	return corsInstance
}

// NewDefaultCORSConfig 初始化并返回全局单例 CORSConfig。
// 后续调用直接返回同一实例。
func NewDefaultCORSConfig() *CORSConfig {
	corsOnce.Do(func() {
		corsInstance = &CORSConfig{
			AllowOrigins:     []string{"http://localhost:5173"},
			AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
			AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
			ExposeHeaders:    []string{"Content-Length"},
			AllowCredentials: true,
		}
	})
	return corsInstance
}

// ApplyFormFile 从 INI 文件的 [cors] 节加载配置并合并到当前实例。
// 支持链式调用。
func (c *CORSConfig) ApplyFormFile(path string) (*CORSConfig, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return c, nil
	}
	sections, err := utils.ParseINI(path)
	if err != nil {
		return nil, err
	}
	if kv, ok := sections["cors"]; ok {
		utils.ApplyOriginList(kv, "allow_origins", &c.AllowOrigins)
		utils.ApplyOriginList(kv, "allow_methods", &c.AllowMethods)
		utils.ApplyOriginList(kv, "allow_headers", &c.AllowHeaders)
		utils.ApplyOriginList(kv, "expose_headers", &c.ExposeHeaders)
		utils.ApplyBool(kv, "allow_credentials", &c.AllowCredentials)
	}
	return c, nil
}

// ApplyFromENV 从 WEAUTH_ 前缀的环境变量加载配置并合并到当前实例。
// 支持链式调用。
func (c *CORSConfig) ApplyFromENV() *CORSConfig {
	utils.EnvOriginList("WEAUTH_CORS_ALLOW_ORIGINS", &c.AllowOrigins)
	utils.EnvOriginList("WEAUTH_CORS_ALLOW_METHODS", &c.AllowMethods)
	utils.EnvOriginList("WEAUTH_CORS_ALLOW_HEADERS", &c.AllowHeaders)
	utils.EnvOriginList("WEAUTH_CORS_EXPOSE_HEADERS", &c.ExposeHeaders)
	utils.EnvBool("WEAUTH_CORS_ALLOW_CREDENTIALS", &c.AllowCredentials)
	return c
}

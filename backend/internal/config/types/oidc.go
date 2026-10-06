package types

import (
	"os"
	"sync"

	"github.com/ctipscn/weauth/backend/internal/config/utils"
)

// OIDCConfig OIDC 认证配置
type OIDCConfig struct {
	ClientID     string
	ClientSecret string
	DiscoveryURL string
	RedirectURL  string
	FrontendURL  string
}

var (
	oidcInstance *OIDCConfig
	oidcOnce     sync.Once
)

// OIDC 返回全局单例 OIDCConfig。
// 必须先调用 NewDefaultOIDCConfig 初始化，否则返回 nil。
func OIDC() *OIDCConfig {
	return oidcInstance
}

// ApplyFormFile 从 INI 文件的 [oidc] 节加载配置并合并到当前实例。
// 支持链式调用。
func (c *OIDCConfig) ApplyFormFile(path string) (*OIDCConfig, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return c, nil
	}
	sections, err := utils.ParseINI(path)
	if err != nil {
		return nil, err
	}
	if kv, ok := sections["oidc"]; ok {
		utils.ApplyString(kv, "client_id", &c.ClientID)
		utils.ApplyString(kv, "client_secret", &c.ClientSecret)
		utils.ApplyString(kv, "discovery_url", &c.DiscoveryURL)
		utils.ApplyString(kv, "redirect_url", &c.RedirectURL)
		utils.ApplyString(kv, "frontend_url", &c.FrontendURL)
	}
	return c, nil
}

// ApplyFromENV 从 WEAUTH_ 前缀的环境变量加载配置并合并到当前实例。
// 支持链式调用。
func (c *OIDCConfig) ApplyFromENV() *OIDCConfig {
	utils.EnvString("WEAUTH_OIDC_CLIENT_ID", &c.ClientID)
	utils.EnvString("WEAUTH_OIDC_CLIENT_SECRET", &c.ClientSecret)
	utils.EnvString("WEAUTH_OIDC_DISCOVERY_URL", &c.DiscoveryURL)
	utils.EnvString("WEAUTH_OIDC_REDIRECT_URL", &c.RedirectURL)
	utils.EnvString("WEAUTH_OIDC_FRONTEND_URL", &c.FrontendURL)
	return c
}

// NewDefaultOIDCConfig 初始化并返回全局单例 OIDCConfig。
// 后续调用直接返回同一实例。
func NewDefaultOIDCConfig() *OIDCConfig {
	oidcOnce.Do(func() {
		oidcInstance = &OIDCConfig{
			ClientID:     "",
			ClientSecret: "",
			DiscoveryURL: "",
			RedirectURL:  "",
			FrontendURL:  "http://localhost:5173",
		}
	})
	return oidcInstance
}
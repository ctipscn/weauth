package config

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/ctipscn/weauth/backend/internal/config/types"
)

// AppConfig 应用总配置。
type AppConfig struct {
	Server    *types.ServerConfig
	Database  *types.DatabaseConfig
	JWT       *types.JWTConfig
	OIDC      *types.OIDCConfig
	CORS      *types.CORSConfig
	RateLimit *types.RateLimitConfig
}

var (
	instance *AppConfig
	once     sync.Once
)

// Config 返回全局单例 AppConfig 实例。
// 必须先调用 Load 完成配置加载，否则返回 nil。
func Config() *AppConfig {
	return instance
}

// NewDefaultConfig 返回一个携带硬编码默认值的新 AppConfig 实例。
//
// Deprecated: 不要直接调用，应通过 Load + GetInstance 获取单例。
func NewDefaultConfig() *AppConfig {
	return &AppConfig{
		Server:    types.NewDefaultServerConfig(),
		Database:  types.NewDefaultDatabaseConfig(),
		JWT:       types.NewDefaultJWTConfig(),
		OIDC:      types.NewDefaultOIDCConfig(),
		CORS:      types.NewDefaultCORSConfig(),
		RateLimit: types.NewDefaultRateLimitConfig(),
	}
}

// Load 按优先级链加载配置并写入全局单例：
//  1. 硬编码默认值
//  2. 系统级配置（/etc/weauth/weauth.conf）
//  3. 用户级配置（$HOME/.config/weauth/weauth.conf）
//  4. 应用专用配置（$WEAUTH_HOME/conf/weauth.conf）
//  5. extraPaths（如 --config 参数）
//  6. 环境变量（WEAUTH_*）
//
// 后面的覆盖前面的。只有首次调用会执行加载；后续调用直接返回已有实例。
func Load(extraPaths ...string) *AppConfig {
	once.Do(func() {
		cfg := NewDefaultConfig()

		// 构建路径列表（低→高优先级）
		var paths []string

		// (a) 系统级
		paths = append(paths, "/etc/weauth/weauth.conf")

		// (b) 用户级
		if home := os.Getenv("HOME"); home != "" {
			paths = append(paths, filepath.Join(home, ".config/weauth/weauth.conf"))
		}

		// (c) 应用专用
		if weauthHome := os.Getenv("WEAUTH_HOME"); weauthHome != "" {
			paths = append(paths, filepath.Join(weauthHome, "conf/weauth.conf"))
		}

		// (d) 附加路径（最高文件优先级）
		paths = append(paths, extraPaths...)

		// 依次尝试，文件不存在或出错时将错误附在节名前以表明出错位置
		for _, p := range paths {
			if p == "" {
				continue
			}
			if _, err := os.Stat(p); os.IsNotExist(err) {
				continue
			}
			if _, err := cfg.ApplyFormFile(p); err != nil {
				// 将文件路径附加到错误信息以明确出错位置
				fmt.Fprintf(os.Stderr, "warn: load config from %s: %v\n", p, err)
			}
		}

		// 环境变量覆盖
		cfg.ApplyFromENV()

		instance = cfg
	})
	return instance
}

// ApplyFormFile 从 INI 文件加载各子配置并合并到当前实例。
// 支持链式调用。
func (c *AppConfig) ApplyFormFile(path string) (*AppConfig, error) {
	if _, err := c.Server.ApplyFormFile(path); err != nil {
		return nil, fmt.Errorf("parse config section [server]: %w", err)
	}
	if _, err := c.Database.ApplyFormFile(path); err != nil {
		return nil, fmt.Errorf("parse config section [database]: %w", err)
	}
	if _, err := c.JWT.ApplyFormFile(path); err != nil {
		return nil, fmt.Errorf("parse config section [jwt]: %w", err)
	}
	if _, err := c.OIDC.ApplyFormFile(path); err != nil {
		return nil, fmt.Errorf("parse config section [oidc]: %w", err)
	}
	if _, err := c.CORS.ApplyFormFile(path); err != nil {
		return nil, fmt.Errorf("parse config section [cors]: %w", err)
	}
	if _, err := c.RateLimit.ApplyFormFile(path); err != nil {
		return nil, fmt.Errorf("parse config section [ratelimit]: %w", err)
	}
	return c, nil
}

// ApplyFromENV 从 WEAUTH_ 前缀的环境变量加载各子配置并合并到当前实例。
// 支持链式调用。
func (c *AppConfig) ApplyFromENV() *AppConfig {
	c.Server.ApplyFromENV()
	c.Database.ApplyFromENV()
	c.JWT.ApplyFromENV()
	c.OIDC.ApplyFromENV()
	c.CORS.ApplyFromENV()
	c.RateLimit.ApplyFromENV()
	return c
}

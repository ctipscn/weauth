package types

import (
	"os"
	"sync"

	"github.com/ctipscn/weauth/backend/internal/config/utils"
)

// ServerConfig 服务器配置
type ServerConfig struct {
	Listen  string
	Port    string
	Mode    string
	LogFile string
}

var (
	srvInstance *ServerConfig
	srvOnce     sync.Once
)

// Server 返回全局单例 ServerConfig。
// 必须先调用 NewDefaultServerConfig 初始化，否则返回 nil。
func Server() *ServerConfig {
	return srvInstance
}

// ApplyFormFile 从 INI 文件的 [server] 节加载配置并合并到当前实例。
// 支持链式调用。
func (c *ServerConfig) ApplyFormFile(path string) (*ServerConfig, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return c, nil
	}
	sections, err := utils.ParseINI(path)
	if err != nil {
		return nil, err
	}
	if kv, ok := sections["server"]; ok {
		utils.ApplyString(kv, "listen", &c.Listen)
		utils.ApplyString(kv, "port", &c.Port)
		utils.ApplyString(kv, "mode", &c.Mode)
		utils.ApplyString(kv, "log_file", &c.LogFile)
	}
	return c, nil
}

// ApplyFromENV 从 WEAUTH_ 前缀的环境变量加载配置并合并到当前实例。
// 支持链式调用。
func (c *ServerConfig) ApplyFromENV() *ServerConfig {
	utils.EnvString("WEAUTH_SERVER_LISTEN", &c.Listen)
	utils.EnvString("WEAUTH_SERVER_PORT", &c.Port)
	utils.EnvString("WEAUTH_SERVER_MODE", &c.Mode)
	utils.EnvString("WEAUTH_SERVER_LOG_FILE", &c.LogFile)
	return c
}

// NewDefaultServerConfig 初始化并返回全局单例 ServerConfig。
// 后续调用直接返回同一实例。
func NewDefaultServerConfig() *ServerConfig {
	srvOnce.Do(func() {
		srvInstance = &ServerConfig{
			Listen:  "127.0.0.1",
			Port:    "8080",
			Mode:    "release",
			LogFile: "logs/app.log",
		}
	})
	return srvInstance
}

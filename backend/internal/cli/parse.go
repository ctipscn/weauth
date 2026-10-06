// Package cli 提供命令行参数解析与验证。
package cli

import (
	"errors"
	"flag"
	"os"
)

// Parse 解析命令行参数，验证必选参数组合。
func Parse() (CliArgs, error) {
	var a CliArgs

	// 内置参数
	flag.StringVar(&a.ConfigPath, "config", "", "配置文件路径")
	flag.BoolVar(&a.Init, "init", false, "初始化用户后退出")
	flag.StringVar(&a.Username, "username", "", "用户名（配合 --init 使用）")
	flag.StringVar(&a.Password, "password", "", "密码（配合 --init 使用）")
	flag.StringVar(&a.FrontendTrialURL, "frontend-trial", "", "创建试用用户和站点后退出（参数为前端 URL）")

	// Server 配置覆盖
	a.Flags.ServerListen = stringFlag("server-listen", "服务监听地址（覆盖配置文件）")
	a.Flags.ServerPort = stringFlag("server-port", "服务监听端口（覆盖配置文件）")

	a.Flags.ServerMode = stringFlag(
		"server-mode", 
		"服务运行模式：debug/release/test（覆盖配置文件）")

	a.Flags.ServerLogFile = stringFlag("server-log-file", "日志文件路径（覆盖配置文件）")

	// Database 配置覆盖
	a.Flags.DatabaseHost = stringFlag("database-host", "数据库主机（覆盖配置文件）")
	a.Flags.DatabasePort = stringFlag("database-port", "数据库端口（覆盖配置文件）")

	a.Flags.DatabaseMode = stringFlag(
		"database-db-mode", "数据库运行模式：debug/release/test（覆盖配置文件）")

	a.Flags.DatabaseUser = stringFlag("database-user", "数据库用户（覆盖配置文件）")
	a.Flags.DatabasePassword = stringFlag("database-password", "数据库密码（覆盖配置文件）")
	a.Flags.DatabaseDBName = stringFlag("database-db-name", "数据库名（覆盖配置文件）")

	a.Flags.DatabaseMaxIdleConns = intFlag(
		"database-max-idle-conns", "最大空闲连接数（覆盖配置文件）")

	a.Flags.DatabaseMaxOpenConns = intFlag(
		"database-max-open-conns", "最大打开连接数（覆盖配置文件）")

	a.Flags.DatabaseConnMaxLifetimeMin = intFlag(
		"database-conn-max-lifetime-minute", "连接最大存活分钟（覆盖配置文件）")

	a.Flags.DatabaseConnMaxIdleMin = intFlag(
		"database-conn-max-idle-minute", "连接最大空闲分钟（覆盖配置文件）")

	// JWT 配置覆盖
	a.Flags.JWTSecret = stringFlag("jwt-secret", "JWT 密钥（覆盖配置文件）")
	a.Flags.JWTExpireHours = intFlag("jwt-expire-hours", "JWT 过期小时数（覆盖配置文件）")

	// OIDC 配置覆盖
	a.Flags.OIDCClientID = stringFlag("oidc-client-id", "OIDC 客户端 ID（覆盖配置文件）")
	a.Flags.OIDCClientSecret = stringFlag("oidc-client-secret", "OIDC 客户端密钥（覆盖配置文件）")
	a.Flags.OIDCDiscoveryURL = stringFlag("oidc-discovery-url", "OIDC 发现 URL（覆盖配置文件）")
	a.Flags.OIDCRedirectURL = stringFlag("oidc-redirect-url", "OIDC 回调 URL（覆盖配置文件）")
	a.Flags.OIDCFrontendURL = stringFlag("oidc-frontend-url", "前端地址（覆盖配置文件）")

	// CORS 配置覆盖
	a.Flags.CORSAllowOrigins = stringFlag("cors-allow-origins", "允许的跨域来源（逗号分隔，覆盖配置文件）")

	// RateLimit 配置覆盖
	a.Flags.RateLimitTimeWindowSeconds = intFlag("ratelimit-time-window-seconds", "限流时间窗口（秒，覆盖配置文件）")
	a.Flags.RateLimitRequestLimit = intFlag("ratelimit-request-limit", "窗口内最大请求次数（覆盖配置文件）")

	flag.Parse()

	switch {
	case a.Init:
		if a.Username == "" {
			return CliArgs{}, errors.New("[CLI] --init requires --username")
		}
		if a.Password == "" {
			return CliArgs{}, errors.New("[CLI] --init requires --password")
		}

	case a.FrontendTrialURL != "":
		// 试用模式无需额外参数

	case a.Username != "" || a.Password != "":
		return CliArgs{}, errors.New("[CLI] --username/--password must be used with --init")
	}

	if a.ConfigPath != "" {
		if _, err := os.Stat(a.ConfigPath); err != nil {
			return CliArgs{}, err
		}
	}

	return a, nil
}

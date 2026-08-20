package config

// DatabaseConfig 数据库连接配置。
type DatabaseConfig struct {
	Driver   string
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	Path     string // SQLite 文件路径
}

// DSN 返回 PostgreSQL 连接字符串。
func (c *DatabaseConfig) DSN() string {
	return "host=" + c.Host +
		" user=" + c.User +
		" password=" + c.Password +
		" dbname=" + c.DBName +
		" port=" + c.Port +
		" sslmode=disable TimeZone=Asia/Shanghai"
}

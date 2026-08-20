package database

import "gorm.io/gorm"

// Storer 存储抽象，供服务和中间件使用。暴露底层 *gorm.DB 用于 GORM 查询，并提供生命周期方法（AutoMigrate、Close）。
type Storer interface {
	// DB 返回底层 GORM 数据库实例。
	DB() *gorm.DB

	// AutoMigrate 对给定模型执行自动迁移。
	AutoMigrate(models ...any) error

	// Close 关闭数据库连接。
	Close() error
}

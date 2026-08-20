package migration

import (
	"github.com/ctipscn/weauth/backend/internal/database"
	"github.com/ctipscn/weauth/backend/internal/models"
)

// Migrate 执行数据库自动迁移，先迁移 User 表再迁移其余业务表。
func Migrate(store database.Storer) error {
	if err := store.AutoMigrate(&models.User{}); err != nil {
		return err
	}

	return store.AutoMigrate(
		&models.Site{},
		&models.SiteDomain{},
		&models.IpPolicy{},
		&models.IpRecord{},
		&models.Challenge{},
		&models.VerificationLog{},
	)
}

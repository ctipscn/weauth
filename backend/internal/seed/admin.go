package seed

import (
	"log"

	"github.com/ctipscn/weauth/backend/internal/models"

	"gorm.io/gorm"
)

// SeedAdmin 确保数据库中存在至少一个管理员用户，并将无主站点分配给默认用户。
func SeedAdmin(gormDB *gorm.DB) {
	var userCount int64
	gormDB.Model(&models.User{}).Count(&userCount)
	if userCount == 0 {
		log.Println("Creating default admin user for existing sites...")
		defaultUser := &models.User{
			Username: "admin",
			Email:    "admin@example.com",
			Role:     "admin",
			Enabled:  true,
		}
		defaultUser.SetPassword("admin123")
		if err := gormDB.Create(defaultUser).Error; err != nil {
			log.Printf("Warning: Failed to create default user: %v", err)
		} else {
			log.Printf("Default admin user created (username: admin, password: admin123)")
		}
	}

	// 为现有的站点分配默认用户
	var sites []models.Site
	gormDB.Where("user_id = 0 OR user_id IS NULL").Find(&sites)
	if len(sites) > 0 {
		var defaultUser models.User
		if err := gormDB.First(&defaultUser).Error; err == nil {
			log.Printf("Assigning %d existing sites to default user...", len(sites))
			for _, site := range sites {
				gormDB.Model(&site).Update("user_id", defaultUser.ID)
			}
			log.Println("Site ownership assignment completed")
		}
	}
}
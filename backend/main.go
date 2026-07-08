package main

import (
	"github.com/ctipscn/weauth/backend/config"
	"github.com/ctipscn/weauth/backend/database"
	"github.com/ctipscn/weauth/backend/handlers"
	"github.com/ctipscn/weauth/backend/middleware"
	"github.com/ctipscn/weauth/backend/models"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// 加载配置
	if err := config.Load(); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 连接数据库
	if err := database.Connect(); err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}

	// 自动迁移数据库 - 先迁移 AdminUser
	if err := database.AutoMigrate(&models.AdminUser{}); err != nil {
		log.Fatalf("Failed to migrate AdminUser: %v", err)
	}

	// 确保至少有一个管理员用户用于现有站点
	var userCount int64
	database.DB.Model(&models.AdminUser{}).Count(&userCount)
	if userCount == 0 {
		log.Println("Creating default admin user for existing sites...")
		defaultUser := &models.AdminUser{
			Username: "admin",
			Email:    "admin@example.com",
			Role:     "admin",
			Enabled:  true,
		}
		defaultUser.SetPassword("admin123")
		if err := database.DB.Create(defaultUser).Error; err != nil {
			log.Printf("Warning: Failed to create default user: %v", err)
		} else {
			log.Printf("Default admin user created (username: admin, password: admin123)")
		}
	}

	// 迁移其他表
	if err := database.AutoMigrate(
		&models.Site{},
		&models.SiteDomain{},
		&models.IpPolicy{},
		&models.IpRecord{},
		&models.Challenge{},
		&models.VerificationLog{},
	); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// 为现有的站点分配默认用户
	var sites []models.Site
	database.DB.Where("user_id = 0 OR user_id IS NULL").Find(&sites)
	if len(sites) > 0 {
		var defaultUser models.AdminUser
		if err := database.DB.First(&defaultUser).Error; err == nil {
			log.Printf("Assigning %d existing sites to default user...", len(sites))
			for _, site := range sites {
				database.DB.Model(&site).Update("user_id", defaultUser.ID)
			}
			log.Println("Site ownership assignment completed")
		}
	}

	log.Println("Database migration completed")

	// 设置 Gin 模式
	gin.SetMode(config.AppConfig.Server.Mode)

	// 创建路由
	r := gin.Default()

	// 中间件
	r.Use(middleware.CORS())

	// 公开接口 (前端调用)
	powGroup := r.Group("/api/pow")
	{
		powGroup.POST("/challenge", handlers.CreateChallenge)
		powGroup.POST("/verify", handlers.VerifyChallenge)
		powGroup.POST("/siteverify", handlers.SiteVerify)
	}

	// 认证接口
	authGroup := r.Group("/api/auth")
	{
		authGroup.GET("/oidc/authorize", handlers.OIDCAuthorize)
		authGroup.GET("/oidc/callback", handlers.OIDCCallback)
		authGroup.POST("/login", handlers.Login)
		authGroup.POST("/logout", handlers.Logout)
		authGroup.GET("/me", middleware.JWTAuth(), handlers.GetCurrentUser)
	}

	// 用户接口 (需要登录)
	userGroup := r.Group("/api/user")
	userGroup.Use(middleware.JWTAuth())
	{
		userGroup.GET("/profile", handlers.GetProfile)
		userGroup.PUT("/password", handlers.UpdatePassword)
	}

	// 管理接口 (需要登录)
	adminGroup := r.Group("/api/admin")
	adminGroup.Use(middleware.JWTAuth())
	{
		// 站点管理
		adminGroup.GET("/sites", handlers.ListSites)
		adminGroup.POST("/sites", handlers.CreateSite)
		adminGroup.GET("/sites/:id", handlers.GetSite)
		adminGroup.PUT("/sites/:id", handlers.UpdateSite)
		adminGroup.DELETE("/sites/:id", handlers.DeleteSite)

		// 域名管理
		adminGroup.GET("/sites/:id/domains", handlers.ListDomains)
		adminGroup.POST("/sites/:id/domains", handlers.AddDomain)
		adminGroup.DELETE("/sites/:id/domains/:domainId", handlers.DeleteDomain)

		// Secret 管理
		adminGroup.GET("/sites/:id/secret", handlers.GetSiteSecret)
		adminGroup.POST("/sites/:id/regenerate-secret", handlers.RegenerateSecret)

		// IP 策略
		adminGroup.GET("/sites/:id/ip-policy", handlers.GetIpPolicy)
		adminGroup.PUT("/sites/:id/ip-policy", handlers.UpdateIpPolicy)

		// IP 记录
		adminGroup.GET("/sites/:id/ip-records", handlers.ListIpRecords)
		adminGroup.DELETE("/sites/:id/ip-records/:ip", handlers.DeleteIpRecord)

		// 统计
		adminGroup.GET("/sites/:id/stats", handlers.GetSiteStats)
	}

	// 启动服务器
	addr := ":" + config.AppConfig.Server.Port
	log.Printf("Server starting on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

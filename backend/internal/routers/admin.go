package routers

import (
	"github.com/ctipscn/weauth/backend/internal/handlers/admin"
	"github.com/ctipscn/weauth/backend/internal/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AdminRegister 注册管理后台路由。
func AdminRegister(r *gin.Engine, ctrl *admin.AdminController, gormDB *gorm.DB) {
	adminGroup := r.Group("/api/admin")
	adminGroup.Use(middleware.JWTAuth(gormDB))
	{
		// 站点管理
		adminGroup.GET("/sites", ctrl.ListSites)
		adminGroup.POST("/sites", ctrl.CreateSite)
		adminGroup.GET("/sites/:id", ctrl.GetSite)
		adminGroup.PUT("/sites/:id", ctrl.UpdateSite)
		adminGroup.DELETE("/sites/:id", ctrl.DeleteSite)

		// 域名管理
		adminGroup.GET("/sites/:id/domains", ctrl.ListDomains)
		adminGroup.POST("/sites/:id/domains", ctrl.AddDomain)
		adminGroup.DELETE("/sites/:id/domains/:domainId", ctrl.DeleteDomain)

		// Secret 管理
		adminGroup.GET("/sites/:id/secret", ctrl.GetSiteSecret)
		adminGroup.POST("/sites/:id/regenerate-secret", ctrl.RegenerateSecret)

		// IP 策略
		adminGroup.GET("/sites/:id/ip-policy", ctrl.GetIpPolicy)
		adminGroup.PUT("/sites/:id/ip-policy", ctrl.UpdateIpPolicy)

		// IP 记录
		adminGroup.GET("/sites/:id/ip-records", ctrl.ListIpRecords)
		adminGroup.DELETE("/sites/:id/ip-records/:ip", ctrl.DeleteIpRecord)

		// 统计
		adminGroup.GET("/sites/:id/stats", ctrl.GetSiteStats)
	}
}
package groups

import (
	consolehandler "github.com/ctipscn/weauth/backend/internal/interface/http/handler/console"
	"github.com/ctipscn/weauth/backend/internal/interface/http/middleware"
	"github.com/gin-gonic/gin"
)

// ConsoleRouteRegister 注册控制台接口（需要登录）
func ConsoleRouteRegister(
	r *gin.Engine,
	consoleHandler *consolehandler.ConsoleHandler,
	authMiddleware *middleware.JWTAuthenticator,
) *gin.Engine {

	consoleGroup := r.Group("/api/console")
	consoleGroup.Use(authMiddleware.RequireAuth())
	{
		// 站点管理
		consoleGroup.GET("/sites", consoleHandler.ListSites)
		consoleGroup.POST("/sites", consoleHandler.CreateSite)
		consoleGroup.GET("/sites/:id", consoleHandler.GetSite)
		consoleGroup.PUT("/sites/:id", consoleHandler.UpdateSite)
		consoleGroup.DELETE("/sites/:id", consoleHandler.DeleteSite)

		// 域名管理
		consoleGroup.GET("/sites/:id/domains", consoleHandler.ListDomains)
		consoleGroup.POST("/sites/:id/domains", consoleHandler.AddDomain)
		consoleGroup.DELETE("/sites/:id/domains/:domainId", consoleHandler.DeleteDomain)

		// Secret 管理
		consoleGroup.GET("/sites/:id/secret", consoleHandler.GetSiteSecret)
		consoleGroup.POST("/sites/:id/regenerate-secret", consoleHandler.RegenerateSecret)

		// IP 策略
		consoleGroup.GET("/sites/:id/ip-policy", consoleHandler.GetIpPolicy)
		consoleGroup.PUT("/sites/:id/ip-policy", consoleHandler.UpdateIpPolicy)

		// IP 记录
		consoleGroup.GET("/sites/:id/ip-records", consoleHandler.ListIpRecords)
		consoleGroup.DELETE("/sites/:id/ip-records/:ip", consoleHandler.DeleteIpRecord)

		// 统计
		consoleGroup.GET("/sites/:id/stats", consoleHandler.GetSiteStats)
	}
	return r
}

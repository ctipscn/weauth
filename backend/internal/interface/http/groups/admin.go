package groups

import (
	adminhandler "github.com/ctipscn/weauth/backend/internal/interface/http/handler/admin"
	"github.com/ctipscn/weauth/backend/internal/interface/http/middleware"
	"github.com/gin-gonic/gin"
)

// AdminRouteRegister 注册管理员后台接口（需要管理员权限）
func AdminRouteRegister(
	r *gin.Engine,
	adminHandler *adminhandler.AdminHandler,
	authMiddleware *middleware.JWTAuthenticator,
) *gin.Engine {

	adminGroup := r.Group("/api/admin")
	adminGroup.Use(authMiddleware.RequireAdmin())
	{
		// 用户管理
		adminGroup.GET("/users", adminHandler.ListUsers)
		adminGroup.DELETE("/users/:id", adminHandler.DeleteUser)

		// 站点管理
		adminGroup.GET("/sites", adminHandler.ListSites)
		adminGroup.GET("/sites/:id", adminHandler.GetSite)
		adminGroup.DELETE("/sites/:id", adminHandler.DeleteSite)
		adminGroup.DELETE("/sites/:id/domains/:domainId", adminHandler.DeleteDomain)
	}
	return r
}

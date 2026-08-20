package routers

import (
	"github.com/ctipscn/weauth/backend/internal/handlers/auth"
	"github.com/ctipscn/weauth/backend/internal/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AuthRegister 注册认证与用户路由。
func AuthRegister(r *gin.Engine, ctrl *auth.AuthController, gormDB *gorm.DB) {
	authGroup := r.Group("/api/auth")
	{
		authGroup.GET("/oidc/authorize", ctrl.OIDCAuthorize)
		authGroup.GET("/oidc/callback", ctrl.OIDCCallback)
		authGroup.POST("/login", ctrl.Login)
		authGroup.POST("/logout", ctrl.Logout)
		authGroup.GET("/me", middleware.JWTAuth(gormDB), ctrl.GetCurrentUser)
	}

	userGroup := r.Group("/api/user")
	userGroup.Use(middleware.JWTAuth(gormDB))
	{
		userGroup.GET("/profile", ctrl.GetProfile)
		userGroup.PUT("/password", ctrl.UpdatePassword)
	}
}
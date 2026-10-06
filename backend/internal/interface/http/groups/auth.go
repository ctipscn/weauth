package groups

import (
	authhandler "github.com/ctipscn/weauth/backend/internal/interface/http/handler/auth"
	"github.com/ctipscn/weauth/backend/internal/interface/http/middleware"
	"github.com/gin-gonic/gin"
)

// AuthRouteRegister 注册认证相关的接口（含 OIDC、登录登出和个人信息查询）
func AuthRouteRegister(
	r *gin.Engine,
	authHandler *authhandler.AuthHandler,
	authMiddleware *middleware.JWTAuthenticator,
) *gin.Engine {

	authGroup := r.Group("/api/auth")
	{
		// authGroup.POST("/register", authHandler.Register) // 邮箱注册暂未开放
		// authGroup.POST("/login", authHandler.Login) // 密码登录暂未开放
		authGroup.POST("/logout", authHandler.Logout)
		authGroup.GET("/oidc/authorize", authHandler.OIDCAuthorize)
		authGroup.GET("/oidc/callback", authHandler.OIDCCallback)
		authGroup.GET("/me", authMiddleware.RequireAuth(), authHandler.GetCurrentUser)
		authGroup.PUT("/password/update", authMiddleware.RequireAuth(), authHandler.UpdatePassword)
		authGroup.PUT("/password/set", authMiddleware.RequireAuth(), authHandler.SetPassword)
	}
	return r
}

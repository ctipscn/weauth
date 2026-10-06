package groups

import (
	userhandler "github.com/ctipscn/weauth/backend/internal/interface/http/handler/user"
	"github.com/ctipscn/weauth/backend/internal/interface/http/middleware"
	"github.com/gin-gonic/gin"
)

// UserRouteRegister 注册用户相关的接口（需要登录）
func UserRouteRegister(
	r *gin.Engine,
	userHandler *userhandler.UserHandler,
	authMiddleware *middleware.JWTAuthenticator,
) *gin.Engine {

	userGroup := r.Group("/api/user")
	userGroup.Use(authMiddleware.RequireAuth())
	{
		userGroup.GET("/profile", userHandler.GetProfile)
		userGroup.PUT("/profile", userHandler.UpdateProfile)
		userGroup.POST("/deactivate", userHandler.DeactivateAccount)
	}
	return r
}

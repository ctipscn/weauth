package http

import (
	"github.com/ctipscn/weauth/backend/internal/interface/http/groups"
	adminhandler "github.com/ctipscn/weauth/backend/internal/interface/http/handler/admin"
	authhandler "github.com/ctipscn/weauth/backend/internal/interface/http/handler/auth"
	consolehandler "github.com/ctipscn/weauth/backend/internal/interface/http/handler/console"
	powhandler "github.com/ctipscn/weauth/backend/internal/interface/http/handler/pow"
	siteverifyhandler "github.com/ctipscn/weauth/backend/internal/interface/http/handler/siteverify"
	userhandler "github.com/ctipscn/weauth/backend/internal/interface/http/handler/user"
	"github.com/ctipscn/weauth/backend/internal/interface/http/middleware"
	"github.com/gin-gonic/gin"
)

// RouteRegister 统一注册所有业务路由
// 各子业务的路由注册分散在 internal/interface/http/groups 子包中，
// 由本函数按顺序统一调用。
func RouteRegister(
	r *gin.Engine,
	authHandler *authhandler.AuthHandler,
	powHandler *powhandler.PoWHandler,
	siteVerifyHandler *siteverifyhandler.SiteVerifyHandler,
	consoleHandler *consolehandler.ConsoleHandler,
	userHandler *userhandler.UserHandler,
	adminHandler *adminhandler.AdminHandler,
	authMiddleware *middleware.JWTAuthenticator,
	trafficLimiter *middleware.TrafficLimiter,
) *gin.Engine {
	r = groups.PoWRouteRegister(r, powHandler, siteVerifyHandler, trafficLimiter)
	r = groups.AuthRouteRegister(r, authHandler, authMiddleware)
	r = groups.UserRouteRegister(r, userHandler, authMiddleware)
	r = groups.ConsoleRouteRegister(r, consoleHandler, authMiddleware)
	r = groups.AdminRouteRegister(r, adminHandler, authMiddleware)
	return r
}

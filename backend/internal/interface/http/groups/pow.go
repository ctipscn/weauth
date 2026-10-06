package groups

import (
	powhandler "github.com/ctipscn/weauth/backend/internal/interface/http/handler/pow"
	siteverifyhandler "github.com/ctipscn/weauth/backend/internal/interface/http/handler/siteverify"
	"github.com/ctipscn/weauth/backend/internal/interface/http/middleware"
	"github.com/gin-gonic/gin"
)

// PoWRouteRegister 注册 PoW（工作量证明）相关的公开接口
func PoWRouteRegister(
	r *gin.Engine, 
	powHandler *powhandler.PoWHandler, 
	siteVerifyHandler *siteverifyhandler.SiteVerifyHandler,
	trafficLimiter *middleware.TrafficLimiter,
) *gin.Engine {
		
	powGroup := r.Group("/api/pow")
	powGroup.Use(trafficLimiter.Handler())
	{
		powGroup.POST("/challenge", powHandler.CreateChallenge)
		powGroup.POST("/verify", powHandler.VerifyChallenge)
		powGroup.POST("/siteverify", siteVerifyHandler.SiteVerify)
	}
	return r
}
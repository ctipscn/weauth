package routers

import (
	"github.com/ctipscn/weauth/backend/internal/handlers/pow"
	"github.com/ctipscn/weauth/backend/internal/handlers/siteverify"

	"github.com/gin-gonic/gin"
)

// PowRegister 注册 PoW 与站点验证路由。
func PowRegister(r *gin.Engine, powCtrl *pow.PowController, siteVerifyCtrl *siteverify.SiteVerifyController) {
	powGroup := r.Group("/api/pow")
	{
		powGroup.POST("/challenge", powCtrl.CreateChallenge)
		powGroup.POST("/verify", powCtrl.VerifyChallenge)
		powGroup.POST("/siteverify", siteVerifyCtrl.SiteVerify)
	}
}
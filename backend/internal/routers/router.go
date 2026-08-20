package routers

import (
	"github.com/ctipscn/weauth/backend/internal/handlers/admin"
	"github.com/ctipscn/weauth/backend/internal/handlers/auth"
	"github.com/ctipscn/weauth/backend/internal/handlers/pow"
	"github.com/ctipscn/weauth/backend/internal/handlers/siteverify"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RouterRegister 依次调用所有子路由注册函数完成路由注册。
func RouterRegister(
	r *gin.Engine,
	authCtrl *auth.AuthController,
	powCtrl *pow.PowController,
	siteVerifyCtrl *siteverify.SiteVerifyController,
	adminCtrl *admin.AdminController,
	gormDB *gorm.DB,
) {
	AuthRegister(r, authCtrl, gormDB)
	PowRegister(r, powCtrl, siteVerifyCtrl)
	AdminRegister(r, adminCtrl, gormDB)
}
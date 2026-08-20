package auth

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"net/url"

	"github.com/ctipscn/weauth/backend/internal/config"
	"github.com/ctipscn/weauth/backend/internal/middleware"
	svcAuth "github.com/ctipscn/weauth/backend/internal/services/auth"

	"github.com/gin-gonic/gin"
)

// AuthController 认证控制器，封装 OIDC 登录、密码登录、用户信息等接口。
// 所有依赖通过构造函数注入，不依赖全局变量。
type AuthController struct {
	authSvc *svcAuth.AuthService
}

// NewAuthController 创建 AuthController 实例。
func NewAuthController(authSvc *svcAuth.AuthService) *AuthController {
	return &AuthController{authSvc: authSvc}
}

// OIDCAuthorize 重定向到 OIDC 授权页面
func (c *AuthController) OIDCAuthorize(ctx *gin.Context) {
	stateBytes := make([]byte, 16)
	rand.Read(stateBytes)
	state := hex.EncodeToString(stateBytes)

	ctx.SetCookie("oauth_state", state, 600, "/", "", false, true)

	authURL, err := c.authSvc.GetOIDCAuthURL(state)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get OIDC auth URL"})
		return
	}

	ctx.Redirect(http.StatusTemporaryRedirect, authURL)
}

// OIDCCallback OIDC 回调处理
func (c *AuthController) OIDCCallback(ctx *gin.Context) {
	code := ctx.Query("code")
	state := ctx.Query("state")
	errorParam := ctx.Query("error")

	frontendURL := config.AppConfig.OIDC.FrontendURL

	if errorParam != "" {
		errorDesc := ctx.Query("error_description")
		redirectURL := frontendURL + "/login?error=" + url.QueryEscape(errorDesc)
		ctx.Redirect(http.StatusTemporaryRedirect, redirectURL)
		return
	}

	savedState, err := ctx.Cookie("oauth_state")
	if err != nil || savedState != state {
		redirectURL := frontendURL + "/login?error=" + url.QueryEscape("Invalid state")
		ctx.Redirect(http.StatusTemporaryRedirect, redirectURL)
		return
	}

	ctx.SetCookie("oauth_state", "", -1, "/", "", false, true)

	tokenResp, err := c.authSvc.ExchangeCodeForToken(code)
	if err != nil {
		redirectURL := frontendURL + "/login?error=" + url.QueryEscape("Failed to exchange code")
		ctx.Redirect(http.StatusTemporaryRedirect, redirectURL)
		return
	}

	userInfo, err := c.authSvc.GetUserInfo(tokenResp.AccessToken)
	if err != nil {
		redirectURL := frontendURL + "/login?error=" + url.QueryEscape("Failed to get user info")
		ctx.Redirect(http.StatusTemporaryRedirect, redirectURL)
		return
	}

	user, err := c.authSvc.FindOrCreateUserByOIDC(userInfo)
	if err != nil {
		redirectURL := frontendURL + "/login?error=" + url.QueryEscape(err.Error())
		ctx.Redirect(http.StatusTemporaryRedirect, redirectURL)
		return
	}

	token, err := middleware.GenerateToken(user.ID, user.Role)
	if err != nil {
		redirectURL := frontendURL + "/login?error=" + url.QueryEscape("Failed to generate token")
		ctx.Redirect(http.StatusTemporaryRedirect, redirectURL)
		return
	}

	ctx.SetCookie("token", token, config.AppConfig.JWT.ExpireHours*3600, "/", "", false, true)
	ctx.Redirect(http.StatusTemporaryRedirect, frontendURL+"/dashboard?token="+token)
}

// Login 密码登录
func (c *AuthController) Login(ctx *gin.Context) {
	var req LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	user, err := c.authSvc.LoginByPassword(req.UserID, req.Password)
	if err != nil {
		status := http.StatusUnauthorized
		message := "Invalid credentials"
		switch err {
		case svcAuth.ErrUserNotFound:
			message = "User not found"
		case svcAuth.ErrUserDisabled:
			status = http.StatusForbidden
			message = "User is disabled"
		case svcAuth.ErrNoPasswordSet:
			message = "Password not set. Please login with E时代通行证 first."
		case svcAuth.ErrInvalidPassword:
			message = "Invalid password"
		}
		ctx.JSON(status, gin.H{"error": message})
		return
	}

	token, err := middleware.GenerateToken(user.ID, user.Role)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	ctx.JSON(http.StatusOK, LoginResponse{
		Token: token,
		User: UserResponse{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			Avatar:   user.Avatar,
			Role:     user.Role,
		},
	})
}

// Logout 登出
func (c *AuthController) Logout(ctx *gin.Context) {
	ctx.SetCookie("token", "", -1, "/", "", false, true)
	ctx.JSON(http.StatusOK, gin.H{"message": "Logged out"})
}

// GetCurrentUser 获取当前用户
func (c *AuthController) GetCurrentUser(ctx *gin.Context) {
	user := middleware.GetCurrentUser(ctx)
	if user == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}

	ctx.JSON(http.StatusOK, UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Avatar:   user.Avatar,
		Role:     user.Role,
	})
}

// GetProfile 获取用户资料
func (c *AuthController) GetProfile(ctx *gin.Context) {
	user := middleware.GetCurrentUser(ctx)
	if user == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id":           user.ID,
		"username":     user.Username,
		"email":        user.Email,
		"avatar":       user.Avatar,
		"role":         user.Role,
		"has_password": user.HasPassword(),
		"created_at":   user.CreatedAt,
	})
}

// UpdatePassword 更新密码
func (c *AuthController) UpdatePassword(ctx *gin.Context) {
	user := middleware.GetCurrentUser(ctx)
	if user == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}

	var req UpdatePasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Password must be at least 6 characters"})
		return
	}

	if err := c.authSvc.UpdateUserPassword(user.ID, req.Password); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Password updated"})
}
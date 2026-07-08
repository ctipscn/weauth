package handlers

import (
	"github.com/ctipscn/weauth/backend/config"
	"github.com/ctipscn/weauth/backend/middleware"
	"github.com/ctipscn/weauth/backend/services"
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

// OIDCAuthorize 重定向到 OIDC 授权页面
func OIDCAuthorize(c *gin.Context) {
	// 生成随机 state
	stateBytes := make([]byte, 16)
	rand.Read(stateBytes)
	state := hex.EncodeToString(stateBytes)

	// 存储 state 到 cookie (用于回调时验证)
	c.SetCookie("oauth_state", state, 600, "/", "", false, true)

	// 获取授权 URL
	authURL, err := services.GetOIDCAuthURL(state)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get OIDC auth URL"})
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, authURL)
}

// OIDCCallback OIDC 回调处理
func OIDCCallback(c *gin.Context) {
	code := c.Query("code")
	state := c.Query("state")
	errorParam := c.Query("error")

	frontendURL := config.AppConfig.OIDC.FrontendURL

	// 检查错误
	if errorParam != "" {
		errorDesc := c.Query("error_description")
		redirectURL := frontendURL + "/login?error=" + url.QueryEscape(errorDesc)
		c.Redirect(http.StatusTemporaryRedirect, redirectURL)
		return
	}

	// 验证 state
	savedState, err := c.Cookie("oauth_state")
	if err != nil || savedState != state {
		redirectURL := frontendURL + "/login?error=" + url.QueryEscape("Invalid state")
		c.Redirect(http.StatusTemporaryRedirect, redirectURL)
		return
	}

	// 清除 state cookie
	c.SetCookie("oauth_state", "", -1, "/", "", false, true)

	// 用 code 换取 token
	tokenResp, err := services.ExchangeCodeForToken(code)
	if err != nil {
		redirectURL := frontendURL + "/login?error=" + url.QueryEscape("Failed to exchange code")
		c.Redirect(http.StatusTemporaryRedirect, redirectURL)
		return
	}

	// 获取用户信息
	userInfo, err := services.GetUserInfo(tokenResp.AccessToken)
	if err != nil {
		redirectURL := frontendURL + "/login?error=" + url.QueryEscape("Failed to get user info")
		c.Redirect(http.StatusTemporaryRedirect, redirectURL)
		return
	}

	// 查找或创建用户
	user, err := services.FindOrCreateUserByOIDC(userInfo)
	if err != nil {
		redirectURL := frontendURL + "/login?error=" + url.QueryEscape(err.Error())
		c.Redirect(http.StatusTemporaryRedirect, redirectURL)
		return
	}

	// 生成 JWT
	token, err := middleware.GenerateToken(user.ID, user.Role)
	if err != nil {
		redirectURL := frontendURL + "/login?error=" + url.QueryEscape("Failed to generate token")
		c.Redirect(http.StatusTemporaryRedirect, redirectURL)
		return
	}

	// 设置 cookie 并重定向到前端
	c.SetCookie("token", token, config.AppConfig.JWT.ExpireHours*3600, "/", "", false, true)
	c.Redirect(http.StatusTemporaryRedirect, frontendURL+"/dashboard?token="+token)
}

// LoginRequest 密码登录请求
type LoginRequest struct {
	UserID   uint   `json:"user_id" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token string `json:"token"`
	User  UserResponse `json:"user"`
}

// UserResponse 用户响应
type UserResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar"`
	Role     string `json:"role"`
}

// Login 密码登录
func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	user, err := services.LoginByPassword(req.UserID, req.Password)
	if err != nil {
		status := http.StatusUnauthorized
		message := "Invalid credentials"
		switch err {
		case services.ErrUserNotFound:
			message = "User not found"
		case services.ErrUserDisabled:
			status = http.StatusForbidden
			message = "User is disabled"
		case services.ErrNoPasswordSet:
			message = "Password not set. Please login with E时代通行证 first."
		case services.ErrInvalidPassword:
			message = "Invalid password"
		}
		c.JSON(status, gin.H{"error": message})
		return
	}

	token, err := middleware.GenerateToken(user.ID, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, LoginResponse{
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
func Logout(c *gin.Context) {
	c.SetCookie("token", "", -1, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "Logged out"})
}

// GetCurrentUser 获取当前用户
func GetCurrentUser(c *gin.Context) {
	user := middleware.GetCurrentUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}

	c.JSON(http.StatusOK, UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Avatar:   user.Avatar,
		Role:     user.Role,
	})
}

// GetProfile 获取用户资料
func GetProfile(c *gin.Context) {
	user := middleware.GetCurrentUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":           user.ID,
		"username":     user.Username,
		"email":        user.Email,
		"avatar":       user.Avatar,
		"role":         user.Role,
		"has_password": user.HasPassword(),
		"created_at":   user.CreatedAt,
	})
}

// UpdatePasswordRequest 更新密码请求
type UpdatePasswordRequest struct {
	Password string `json:"password" binding:"required,min=6"`
}

// UpdatePassword 更新密码
func UpdatePassword(c *gin.Context) {
	user := middleware.GetCurrentUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}

	var req UpdatePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password must be at least 6 characters"})
		return
	}

	if err := services.UpdateUserPassword(user.ID, req.Password); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password updated"})
}

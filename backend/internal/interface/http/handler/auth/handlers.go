package auth

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"net/url"

	"github.com/ctipscn/weauth/backend/internal/application/auth"
	"github.com/ctipscn/weauth/backend/internal/config"
	infralogger "github.com/ctipscn/weauth/backend/internal/infrastructure/logger"
	"github.com/ctipscn/weauth/backend/internal/interface/http/middleware"
	"github.com/ctipscn/weauth/backend/internal/interface/http/request"
	"github.com/ctipscn/weauth/backend/internal/interface/http/response"

	"github.com/gin-gonic/gin"
)

// AuthHandler 认证 HTTP 处理器
type AuthHandler struct {
	authApp        *auth.AuthApp
	authMiddleware *middleware.JWTAuthenticator
	log            *infralogger.LogService
}

// NewAuthHandler 创建认证 HTTP 处理器
func NewAuthHandler(authApp *auth.AuthApp, authMiddleware *middleware.JWTAuthenticator, log *infralogger.LogService) *AuthHandler {
	return &AuthHandler{authApp: authApp, authMiddleware: authMiddleware, log: log}
}

// OIDCAuthorize 重定向到 OIDC 授权页面
func (h *AuthHandler) OIDCAuthorize(c *gin.Context) {
	h.log.Infof("[HTTP-AUTH] OIDCAuthorize entry")
	// 生成随机 state
	stateBytes := make([]byte, 16)
	rand.Read(stateBytes)
	state := hex.EncodeToString(stateBytes)

	// 存储 state 到 cookie (用于回调时验证)
	c.SetCookie("oauth_state", state, 600, "/", "", false, true)

	// 获取授权 URL
	oidcCfg := config.Config().OIDC
	authURL, err := h.authApp.OIDCClient().GetAuthURL(oidcCfg.ClientID, oidcCfg.RedirectURL, state)
	if err != nil {
		h.log.Errorf("[HTTP-AUTH] OIDCAuthorize: failed to get auth URL: %v", err)
		c.JSON(
			http.StatusInternalServerError,
			response.ErrorBaseResponse(0, 0, "Failed to get OIDC auth URL"),
		)
		return
	}

	h.log.Infof("[HTTP-AUTH] OIDCAuthorize: redirecting to %s", authURL)
	c.Redirect(http.StatusTemporaryRedirect, authURL)
}

// OIDCCallback OIDC 回调处理
func (h *AuthHandler) OIDCCallback(c *gin.Context) {
	code := c.Query("code")
	state := c.Query("state")
	errorParam := c.Query("error")

	h.log.Infof("[HTTP-AUTH] OIDCCallback entry (code=%s..., state=%s...)", safeTruncate(code, 8), safeTruncate(state, 8))

	frontendURL := config.Config().OIDC.FrontendURL

	// 检查错误
	if errorParam != "" {
		errorDesc := c.Query("error_description")
		h.log.Warnf("[HTTP-AUTH] OIDCCallback: provider error (error=%s, description=%s)", errorParam, errorDesc)
		redirectURL := frontendURL + "/login?error=" + url.QueryEscape(errorDesc)
		c.Redirect(http.StatusTemporaryRedirect, redirectURL)
		return
	}

	// 验证 state
	savedState, err := c.Cookie("oauth_state")
	if err != nil || savedState != state {
		h.log.Warnf("[HTTP-AUTH] OIDCCallback: state mismatch")
		redirectURL := frontendURL + "/login?error=" + url.QueryEscape("Invalid state")
		c.Redirect(http.StatusTemporaryRedirect, redirectURL)
		return
	}

	// 清除 state cookie
	c.SetCookie("oauth_state", "", -1, "/", "", false, true)

	// 用 code 换取 token
	oidcCfg := config.Config().OIDC
	tokenResp, err := h.authApp.OIDCClient().ExchangeCode(
		oidcCfg.ClientID,
		oidcCfg.ClientSecret,
		oidcCfg.RedirectURL,
		code,
	)

	if err != nil {
		h.log.Errorf("[HTTP-AUTH] OIDCCallback: failed to exchange code: %v", err)
		redirectURL := frontendURL + "/login?error=" + url.QueryEscape("Failed to exchange code")
		c.Redirect(http.StatusTemporaryRedirect, redirectURL)
		return
	}

	// 获取用户信息
	userInfo, err := h.authApp.OIDCClient().GetUserInfo(tokenResp.AccessToken)
	if err != nil {
		h.log.Errorf("[HTTP-AUTH] OIDCCallback: failed to get user info: %v", err)
		redirectURL := frontendURL + "/login?error=" + url.QueryEscape("Failed to get user info")
		c.Redirect(http.StatusTemporaryRedirect, redirectURL)
		return
	}

	// 查找或创建用户
	user, err := h.authApp.FindOrCreateUserByOIDC(userInfo)
	if err != nil {
		h.log.Warnf("[HTTP-AUTH] OIDCCallback: user lookup/creation failed: %v", err)
		redirectURL := frontendURL + "/login?error=" + url.QueryEscape(err.Error())
		c.Redirect(http.StatusTemporaryRedirect, redirectURL)
		return
	}

	// 生成 JWT
	token, err := h.authMiddleware.GenerateToken(user.ID, user.Role.String())
	if err != nil {
		h.log.Errorf("[HTTP-AUTH] OIDCCallback: failed to generate token: %v", err)
		redirectURL := frontendURL + "/login?error=" + url.QueryEscape("Failed to generate token")
		c.Redirect(http.StatusTemporaryRedirect, redirectURL)
		return
	}

	h.log.Infof("[HTTP-AUTH] OIDCCallback: success (userID=%d)", user.ID)
	// 设置 cookie 并重定向到前端
	c.SetCookie("token", token, config.Config().JWT.ExpireHours*3600, "/", "", false, true)
	c.Redirect(http.StatusTemporaryRedirect, frontendURL+"/auth/callback?token="+token)
}

// Register 邮箱注册（用户名默认设置为邮箱）
func (h *AuthHandler) Register(c *gin.Context) {
	_, req, err := request.BindJSON[RegisterRequest](c)
	if err != nil {
		h.log.Warnf("[HTTP-AUTH] Register: invalid request body")
		c.JSON(http.StatusBadRequest, response.ErrorBaseResponse(0, 0, "Invalid request"))
		return
	}
	h.log.Infof("[HTTP-AUTH] Register entry (email=%s)", req.Email)

	user, err := h.authApp.Register(req.Email, req.Email, req.Password)
	if err != nil {
		status := http.StatusBadRequest
		message := err.Error()
		switch err {
		case auth.ErrEmailTaken:
			message = "Email already taken"
		case auth.ErrInvalidEmail:
			message = "Invalid email address"
		case auth.ErrPasswordTooShort:
			message = "Password must be at least 6 characters"
		}
		h.log.Warnf("[HTTP-AUTH] Register failed: %s", message)
		c.JSON(status, response.ErrorBaseResponse(0, 0, message))
		return
	}

	token, err := h.authMiddleware.GenerateToken(user.ID, user.Role.String())
	if err != nil {
		h.log.Errorf("[HTTP-AUTH] Register: failed to generate token (userID=%d): %v", user.ID, err)
		c.JSON(http.StatusInternalServerError, response.ErrorBaseResponse(0, 0, "Failed to generate token"))
		return
	}

	h.log.Infof("[HTTP-AUTH] Register success (userID=%d, email=%s)", user.ID, req.Email)
	c.JSON(http.StatusCreated, response.SuccessBaseResponse(RegisterResponse{
		Token: token,
		User: UserResponse{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			Avatar:   user.Avatar,
			Role:     user.Role.String(),
		},
	}))
}

// Login 邮箱密码登录
func (h *AuthHandler) Login(c *gin.Context) {
	_, req, err := request.BindJSON[LoginRequest](c)
	if err != nil {
		h.log.Warnf("[HTTP-AUTH] Login: invalid request body")
		c.JSON(http.StatusBadRequest, response.ErrorBaseResponse(0, 0, "Invalid request"))
		return
	}
	h.log.Infof("[HTTP-AUTH] Login entry (email_or_username=%s)", req.EmailOrUsername)

	user, err := h.authApp.LoginByEmailOrUsername(req.EmailOrUsername, req.Password)
	if err != nil {
		status := http.StatusUnauthorized
		message := "Invalid email/username or password"
		switch err {
		case auth.ErrUserNotFound:
			message = "User not found"
		case auth.ErrUserDisabled:
			status = http.StatusForbidden
			message = "User is disabled"
		case auth.ErrNoPasswordSet:
			message = "Password not set. Please login with OIDC first."
		case auth.ErrInvalidPassword:
			message = "Invalid password"
		}
		h.log.Warnf("[HTTP-AUTH] Login failed: %s (email_or_username=%s)", message, req.EmailOrUsername)
		c.JSON(status, response.ErrorBaseResponse(0, 0, message))
		return
	}

	token, err := h.authMiddleware.GenerateToken(user.ID, user.Role.String())
	if err != nil {
		h.log.Errorf("[HTTP-AUTH] Login: failed to generate token (userID=%d): %v", user.ID, err)
		c.JSON(http.StatusInternalServerError, response.ErrorBaseResponse(0, 0, "Failed to generate token"))
		return
	}

	h.log.Infof("[HTTP-AUTH] Login success (userID=%d, email_or_username=%s)", user.ID, req.EmailOrUsername)
	c.JSON(http.StatusOK, response.SuccessBaseResponse(LoginResponse{
		Token: token,
		User: UserResponse{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			Avatar:   user.Avatar,
			Role:     user.Role.String(),
		},
	}))
}

// Logout 登出
func (h *AuthHandler) Logout(c *gin.Context) {
	user := middleware.GetCurrentUser(c)
	if user != nil {
		h.log.Infof("[HTTP-AUTH] Logout (userID=%d)", user.ID)
	}
	c.SetCookie("token", "", -1, "/", "", false, true)
	c.JSON(http.StatusOK, response.SuccessBaseResponse("Logged out"))
}

// GetCurrentUser 获取当前用户
func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
	user := middleware.GetCurrentUser(c)
	if user == nil {
		h.log.Warnf("[HTTP-AUTH] GetCurrentUser: not authenticated")
		c.JSON(http.StatusUnauthorized, response.ErrorBaseResponse(0, 0, "Not authenticated"))
		return
	}

	h.log.Infof("[HTTP-AUTH] GetCurrentUser success (userID=%d)", user.ID)
	c.JSON(http.StatusOK, response.SuccessBaseResponse(UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Avatar:   user.Avatar,
		Role:     user.Role.String(),
	}))
}

// SetPassword OIDC-only 用户设置密码（无需旧密码）
func (h *AuthHandler) SetPassword(c *gin.Context) {
	user := middleware.GetCurrentUser(c)
	if user == nil {
		h.log.Warnf("[HTTP-AUTH] SetPassword: not authenticated")
		c.JSON(http.StatusUnauthorized, response.ErrorBaseResponse(0, 0, "Not authenticated"))
		return
	}

	_, req, err := request.BindJSON[SetPasswordRequest](c)
	if err != nil {
		h.log.Warnf("[HTTP-AUTH] SetPassword: invalid request (userID=%d)", user.ID)
		c.JSON(http.StatusBadRequest, response.ErrorBaseResponse(0, 0, "New password must be at least 6 characters"))
		return
	}

	if err := h.authApp.SetUserPassword(user.ID, req.NewPassword); err != nil {
		status := http.StatusInternalServerError
		message := "Failed to set password"
		switch err {
		case auth.ErrPasswordAlreadySet:
			status = http.StatusBadRequest
			message = "Password already set or not an OIDC user, use /api/auth/password/update to change password"
		case auth.ErrUserNotFound:
			status = http.StatusNotFound
			message = "User not found"
		}
		h.log.Errorf("[HTTP-AUTH] SetPassword: failed (userID=%d): %v", user.ID, err)
		c.JSON(status, response.ErrorBaseResponse(0, 0, message))
		return
	}

	h.log.Infof("[HTTP-AUTH] Password set (userID=%d)", user.ID)
	c.JSON(http.StatusOK, response.SuccessBaseResponse("Password set"))
}

// UpdatePassword 更新密码（需要验证旧密码）
func (h *AuthHandler) UpdatePassword(c *gin.Context) {
	user := middleware.GetCurrentUser(c)
	if user == nil {
		h.log.Warnf("[HTTP-AUTH] UpdatePassword: not authenticated")
		c.JSON(http.StatusUnauthorized, response.ErrorBaseResponse(0, 0, "Not authenticated"))
		return
	}

	_, req, err := request.BindJSON[UpdatePasswordRequest](c)
	if err != nil {
		h.log.Warnf("[HTTP-AUTH] UpdatePassword: invalid request (userID=%d)", user.ID)
		c.JSON(http.StatusBadRequest, response.ErrorBaseResponse(0, 0, "New password must be at least 6 characters"))
		return
	}

	if err := h.authApp.UpdateUserPassword(user.ID, req.OldPassword, req.NewPassword); err != nil {
		status := http.StatusInternalServerError
		message := "Failed to update password"
		if err == auth.ErrOldPasswordMismatch {
			status = http.StatusBadRequest
			message = "Old password is incorrect"
		}
		h.log.Errorf("[HTTP-AUTH] UpdatePassword: failed (userID=%d): %v", user.ID, err)
		c.JSON(status, response.ErrorBaseResponse(0, 0, message))
		return
	}

	h.log.Infof("[HTTP-AUTH] Password updated (userID=%d)", user.ID)
	c.JSON(http.StatusOK, response.SuccessBaseResponse("Password updated"))
}

// safeTruncate 截断字符串用于日志，避免记录完整 token 或敏感信息
func safeTruncate(s string, maxLen int) string {
	if s == "" {
		return ""
	}
	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}
	return string(runes[:maxLen]) + "..."
}

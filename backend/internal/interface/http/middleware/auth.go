package middleware

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/ctipscn/weauth/backend/internal/domain/auth"
	"github.com/ctipscn/weauth/backend/internal/domain/user"
	infralogger "github.com/ctipscn/weauth/backend/internal/infrastructure/logger"
	"github.com/ctipscn/weauth/backend/internal/interface/http/response"

	"github.com/gin-gonic/gin"
	jwtlib "github.com/golang-jwt/jwt/v5"
)

// Sentinel errors returned by Validate.
var (
	ErrAuthTokenRequired     = errors.New("authorization token required")
	ErrInvalidOrExpiredToken = errors.New("invalid or expired token")
	ErrUserNotFound          = errors.New("user not found")
	ErrUserDisabled          = errors.New("user is disabled")
)

// claims JWT 声明
type claims struct {
	UserID      uint   `json:"user_id"`
	Role        string `json:"role"`
	IsOIDCLogin bool   `json:"is_oidc_login"`
	jwtlib.RegisteredClaims
}

// AuthConfig JWT 认证中间件配置
type AuthConfig struct {
	Secret      string
	ExpireHours int
}

// JWTAuthenticator 认证中间件，提供 Token 校验、用户身份注入和角色鉴权。
// 通过 RequireAuth / RequireAdmin 生成 gin 中间件，验证失败时直接写入 HTTP 响应并终止请求链。
type JWTAuthenticator struct {
	userRepo user.Repository
	authRepo auth.Repository
	log      *infralogger.LogService
	cfg      *AuthConfig
}

// NewJWTAuthenticator 创建 JWTAuthenticator
func NewJWTAuthenticator(
	userRepo user.Repository, 
	authRepo auth.Repository, 
	log *infralogger.LogService,
	cfg *AuthConfig,
) *JWTAuthenticator {
	return &JWTAuthenticator{userRepo: userRepo, authRepo: authRepo, log: log, cfg: cfg}
}

// Validate 校验 JWT token 并返回用户
func (m *JWTAuthenticator) Validate(c *gin.Context) (*user.User, error) {
	tokenString := extractToken(c)
	if tokenString == "" {
		return nil, ErrAuthTokenRequired
	}

	claims, err := m.ParseToken(tokenString)
	if err != nil {
		return nil, ErrInvalidOrExpiredToken
	}

	u, err := m.userRepo.FindByID(claims.UserID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	if !u.Enabled {
		return nil, ErrUserDisabled
	}

	m.log.Infof("[AUTH-MW] Validate success (userID=%d, role=%s)", u.ID, u.Role.String())
	return u, nil
}

// RequireAuth 返回 gin 中间件，要求请求必须携带有效的登录凭证。
// 验证失败时直接写入 HTTP 响应并终止请求链，不继续执行后续 Handler。
func (m *JWTAuthenticator) RequireAuth() gin.HandlerFunc {
	return m.requireAuth(false)
}

// RequireAdmin 返回 gin 中间件，要求请求携带的登录凭证属于管理员角色。
// 非管理员或未登录均直接拦截。
func (m *JWTAuthenticator) RequireAdmin() gin.HandlerFunc {
	return m.requireAuth(true)
}

func (m *JWTAuthenticator) requireAuth(requireAdmin bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		u, err := m.Validate(c)
		if err != nil {
			writeAuthError(c, err)
			c.Abort()
			return
		}

		if requireAdmin && !u.IsAdmin() {
			c.JSON(
				http.StatusForbidden,
				response.ErrorBaseResponse(0, 0, "Admin privileges required"),
			)
			c.Abort()
			return
		}

		c.Set("user", u)
		c.Set("userID", u.ID)
		c.Set("userRole", u.Role.String())
		c.Next()
	}
}

// writeAuthError 将 Validate 返回的 sentinel error 映射为 HTTP 响应。
func writeAuthError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, ErrAuthTokenRequired),
		errors.Is(err, ErrInvalidOrExpiredToken),
		errors.Is(err, ErrUserNotFound):
		c.JSON(http.StatusUnauthorized, response.ErrorBaseResponse(0, 0, "Unauthorized"))
	case errors.Is(err, ErrUserDisabled):
		c.JSON(http.StatusForbidden, response.ErrorBaseResponse(0, 0, "User is disabled"))
	default:
		c.JSON(http.StatusInternalServerError, response.ErrorBaseResponse(0, 0, "Internal error"))
	}
}

// extractToken 从请求中提取 token
func extractToken(c *gin.Context) string {
	// 从 Authorization header 获取
	authHeader := c.GetHeader("Authorization")
	if authHeader != "" {
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
			return parts[1]
		}
	}

	// 从 cookie 获取
	if token, err := c.Cookie("token"); err == nil {
		return token
	}

	// 从查询参数获取
	return c.Query("token")
}

// GenerateToken 生成 JWT token，通过 authRepo 查询 Auth 领域对象确定 is_oidc_login
func (m *JWTAuthenticator) GenerateToken(userID uint, role string) (string, error) {
	isOIDCLogin := false
	cred, err := m.authRepo.FindByUserID(userID)
	if err == nil {
		isOIDCLogin = cred.IsOIDCLogin()
	}

	claims := claims{
		UserID:      userID,
		Role:        role,
		IsOIDCLogin: isOIDCLogin,
		RegisteredClaims: jwtlib.RegisteredClaims{
			ExpiresAt: jwtlib.NewNumericDate(
				time.Now().Add(
					time.Duration(m.cfg.ExpireHours) * time.Hour)),

			IssuedAt: jwtlib.NewNumericDate(time.Now()),
		},
	}

	token := jwtlib.NewWithClaims(jwtlib.SigningMethodHS256, claims)
	return token.SignedString([]byte(m.cfg.Secret))
}

// ParseToken 解析 JWT token
func (m *JWTAuthenticator) ParseToken(tokenString string) (*claims, error) {
	token, err := jwtlib.ParseWithClaims(
		tokenString, &claims{},
		func(token *jwtlib.Token) (any, error) {
			return []byte(m.cfg.Secret), nil
		})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*claims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwtlib.ErrSignatureInvalid
}

// GetCurrentUser 从上下文获取当前用户
func GetCurrentUser(c *gin.Context) *user.User {
	if u, exists := c.Get("user"); exists {
		return u.(*user.User)
	}
	return nil
}

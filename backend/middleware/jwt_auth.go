package middleware

import (
	"github.com/ctipscn/weauth/backend/config"
	"github.com/ctipscn/weauth/backend/database"
	"github.com/ctipscn/weauth/backend/models"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// JWTAuth JWT 鉴权中间件
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := extractToken(c)
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token required"})
			c.Abort()
			return
		}

		claims, err := ParseToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// 查询用户
		var user models.AdminUser
		if err := database.DB.First(&user, claims.UserID).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			c.Abort()
			return
		}

		if !user.Enabled {
			c.JSON(http.StatusForbidden, gin.H{"error": "User is disabled"})
			c.Abort()
			return
		}

		// 设置用户信息到上下文
		c.Set("user", &user)
		c.Set("userID", user.ID)
		c.Set("userRole", user.Role)

		c.Next()
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

// GenerateToken 生成 JWT token
func GenerateToken(userID uint, role string) (string, error) {
	claims := Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(config.AppConfig.JWT.ExpireHours) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.AppConfig.JWT.Secret))
}

// ParseToken 解析 JWT token
func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.AppConfig.JWT.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrSignatureInvalid
}

// GetCurrentUser 从上下文获取当前用户
func GetCurrentUser(c *gin.Context) *models.AdminUser {
	if user, exists := c.Get("user"); exists {
		return user.(*models.AdminUser)
	}
	return nil
}

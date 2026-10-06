package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// CORSConfig 跨域中间件配置
type CORSConfig struct {
	AllowOrigins     []string
	AllowMethods     []string
	AllowHeaders     []string
	ExposeHeaders    []string
	AllowCredentials bool
}

type CORSInterceptor struct {
	cfg *CORSConfig
}

// NewCORSInterceptor 返回跨域拦截器
func NewCORSInterceptor(cfg *CORSConfig) *CORSInterceptor {
	return &CORSInterceptor{cfg: cfg}
}

func (c *CORSInterceptor) ToGinCORS() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:     c.cfg.AllowOrigins,
		AllowMethods:     c.cfg.AllowMethods,
		AllowHeaders:     c.cfg.AllowHeaders,
		ExposeHeaders:    c.cfg.ExposeHeaders,
		AllowCredentials: c.cfg.AllowCredentials,
	})
}

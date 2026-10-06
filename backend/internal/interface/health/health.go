package health

import (
	"net/http"
	"time"

	infrahealth "github.com/ctipscn/weauth/backend/internal/infrastructure/health"
	infralogger "github.com/ctipscn/weauth/backend/internal/infrastructure/logger"
	"github.com/gin-gonic/gin"
)

// DefaultHealthAddr 健康检查端点默认监听地址。
const DefaultHealthAddr = "127.0.0.1:8081"

// RegisterHealthRoute 创建健康检查 HTTP 服务并异步启动。
func RegisterHealthRoute(
	checker *infrahealth.HealthChecker,
	logService *infralogger.LogService,
	startTime time.Time,
	addr string,
) *http.Server {
	if addr == "" {
		addr = DefaultHealthAddr
	}

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	r.GET("/health", func(c *gin.Context) {
		status := checker.Snapshot(startTime)
		if status.Status != "ok" {
			c.JSON(http.StatusServiceUnavailable, status)
			return
		}
		c.JSON(http.StatusOK, status)
	})

	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	go func() {
		logService.Info("[HEALTH] Health check server starting on " + addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logService.Error("[HEALTH] Health check server error: " + err.Error())
		}
	}()

	return srv
}

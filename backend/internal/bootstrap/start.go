package bootstrap

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	appadmin "github.com/ctipscn/weauth/backend/internal/application/admin"
	"github.com/ctipscn/weauth/backend/internal/application/auth"
	appoidc "github.com/ctipscn/weauth/backend/internal/application/oidc"
	"github.com/ctipscn/weauth/backend/internal/application/pow"
	"github.com/ctipscn/weauth/backend/internal/application/site"
	appuser "github.com/ctipscn/weauth/backend/internal/application/user"
	"github.com/ctipscn/weauth/backend/internal/config"
	auth_domain "github.com/ctipscn/weauth/backend/internal/domain/auth"
	domainuser "github.com/ctipscn/weauth/backend/internal/domain/user"
	"github.com/ctipscn/weauth/backend/internal/infrastructure"
	"github.com/ctipscn/weauth/backend/internal/infrastructure/database"
	challengerepo "github.com/ctipscn/weauth/backend/internal/infrastructure/database/repos/challenge"
	ipolicyrepo "github.com/ctipscn/weauth/backend/internal/infrastructure/database/repos/ippolicy"
	irecordrepo "github.com/ctipscn/weauth/backend/internal/infrastructure/database/repos/iprecord"
	siterepo "github.com/ctipscn/weauth/backend/internal/infrastructure/database/repos/site"
	verifylogrepo "github.com/ctipscn/weauth/backend/internal/infrastructure/database/repos/verifylog"
	infrahealth "github.com/ctipscn/weauth/backend/internal/infrastructure/health"
	infralogger "github.com/ctipscn/weauth/backend/internal/infrastructure/logger"
	"github.com/ctipscn/weauth/backend/internal/infrastructure/server"
	"github.com/ctipscn/weauth/backend/internal/infrastructure/security"
	"github.com/ctipscn/weauth/backend/internal/interface/health"
	route "github.com/ctipscn/weauth/backend/internal/interface/http"
	adminhandler "github.com/ctipscn/weauth/backend/internal/interface/http/handler/admin"
	authhandler "github.com/ctipscn/weauth/backend/internal/interface/http/handler/auth"
	consolehandler "github.com/ctipscn/weauth/backend/internal/interface/http/handler/console"
	powhandler "github.com/ctipscn/weauth/backend/internal/interface/http/handler/pow"
	siteverifyhandler "github.com/ctipscn/weauth/backend/internal/interface/http/handler/siteverify"
	userhandler "github.com/ctipscn/weauth/backend/internal/interface/http/handler/user"
	"github.com/ctipscn/weauth/backend/internal/interface/http/middleware"
)

// fatal 负责记录错误日志后终止进程，统一 error → exit 路径。
func fatal(logService *infralogger.LogService, msg string) {
	logService.Error(msg)
	log.Fatal(msg)
}

// StartApp 构建依赖图并启动 HTTP 服务。
func StartApp(
	logService *infralogger.LogService,
	userRepo domainuser.Repository,
	authCredRepo auth_domain.Repository,
	boot *Shutdown,
	dbDriver *database.PostgresDriver,
) {
	// 应用运行模式（所有组件共享）
	appMode, err := infrastructure.StringToMode(config.Config().Server.Mode)
	if err != nil {
		fatal(
			logService, fmt.Sprintf(
				"[BOOT] Invalid server mode %q: %v",
				config.Config().Server.Mode,
				err,
			),
		)
	}

	// 健康检查服务
	healthChecker := infrahealth.NewHealthChecker()
	appStartTime := time.Now()
	srvHealth := health.RegisterHealthRoute(healthChecker, logService, appStartTime, "")
	boot.RegisterCleanupFn(
		"Shutdown health check server",
		func(ctx context.Context) error {
			return srvHealth.Shutdown(ctx)
		},
	)

	// 仓储
	siteRepo := siterepo.NewSiteRepo(dbDriver)
	challengeRepo := challengerepo.NewChallengeRepo(dbDriver)
	ipPolicyRepo := ipolicyrepo.NewIpPolicyRepo(dbDriver)
	ipRecordRepo := irecordrepo.NewIpRecordRepo(dbDriver)
	verifyLogRepo := verifylogrepo.NewVerifyLogRepo(dbDriver)

	// 应用服务
	powApp := pow.NewPoWApp(
		challengeRepo,
		siteRepo,
		ipPolicyRepo,
		ipRecordRepo,
		verifyLogRepo,
		logService,
	)

	siteApp := site.NewSiteApp(siteRepo, ipPolicyRepo, ipRecordRepo, powApp, logService)

	passwdValid := security.NewEncryptor()

	oidcClient := appoidc.NewOIDCClient()
	oidcClient.SetDiscoveryURL(config.Config().OIDC.DiscoveryURL)
	authApp := auth.NewAuthApp(
		userRepo,
		authCredRepo,
		oidcClient,
		config.Config().OIDC.FrontendURL,
		passwdValid,
		logService)

	userApp := appuser.NewUserApp(userRepo, authCredRepo, logService)

	// 管理员应用服务
	adminApp := appadmin.NewAdminApp(
		dbDriver.DB(),
		userRepo,
		siteRepo,
		authCredRepo,
		logService,
	)

	// 中间件 / 拦截器（先于处理器创建，因为 handler 需要引用中间件）
	corsMiddleware := middleware.NewCORSInterceptor(
		&middleware.CORSConfig{
			AllowOrigins:     config.Config().CORS.AllowOrigins,
			AllowHeaders:     config.Config().CORS.AllowHeaders,
			AllowMethods:     config.Config().CORS.AllowMethods,
			ExposeHeaders:    config.Config().CORS.ExposeHeaders,
			AllowCredentials: config.Config().CORS.AllowCredentials,
		},
	)

	authMiddleware := middleware.NewJWTAuthenticator(
		userRepo,
		authCredRepo,
		logService,
		&middleware.AuthConfig{
			Secret:      config.Config().JWT.Secret,
			ExpireHours: config.Config().JWT.ExpireHours,
		},
	)

	trafficLimiter := middleware.NewTrafficLimiter(
		&middleware.LimiterConfig{
			TimeWindowSeconds: config.Config().RateLimit.TimeWindowSeconds,
			RequestLimit:      config.Config().RateLimit.RequestLimit,
		},
	)
	go trafficLimiter.Cleanup(time.Duration(config.Config().RateLimit.TimeWindowSeconds) * time.Second)

	// HTTP 处理器
	authHandler := authhandler.NewAuthHandler(authApp, authMiddleware, logService)
	adminHandler := adminhandler.NewAdminHandler(adminApp, logService)
	consoleHandler := consolehandler.NewConsoleHandler(siteApp, powApp, logService)
	powHandler := powhandler.NewPoWHandler(powApp, logService)
	siteVerifyHandler := siteverifyhandler.NewSiteVerifyHandler(powApp, logService)
	userHandler := userhandler.NewUserHandler(userApp, logService)

	// Gin 服务器 + 路由
	addr := config.Config().Server.Listen + ":" + config.Config().Server.Port
	gs := server.NewGinServer(addr)

	// 依赖注入 CORS 策略
	gs.Engine().Use(corsMiddleware.ToGinCORS())

	// 应用运行模式到 Gin
	if err := gs.ApplyMode(appMode); err != nil {
		fatal(logService, fmt.Sprintf("[BOOT] Failed to apply mode to Gin server: %v", err))
	}
	logService.Info(fmt.Sprintf("[BOOT] Gin server mode set to %s", appMode))

	route.RouteRegister(
		gs.Engine(),
		authHandler,
		powHandler,
		siteVerifyHandler,
		consoleHandler,
		userHandler,
		adminHandler,
		authMiddleware,
		trafficLimiter,
	)

	// 应用运行模式到数据库
	if err := dbDriver.ApplyMode(appMode); err != nil {
		fatal(logService, fmt.Sprintf("[BOOT] Failed to apply mode to database: %v", err))
	}
	logService.Info(fmt.Sprintf("[BOOT] Database mode set to %s", appMode))

	// 注册 HTTP 清理
	boot.RegisterCleanupFn("Shutdown HTTP server", func(ctx context.Context) error {
		return gs.Shutdown(ctx)
	})

	// 创建 Listener，确认端口可用后打印监听地址
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		fatal(logService, fmt.Sprintf("[BOOT] Failed to listen on %s: %v", addr, err))
	}
	fmt.Printf("[BOOT] HTTP server is listening on %s\n", addr)
	logService.Info(fmt.Sprintf("[BOOT] HTTP server is listening on %s", addr))

	// 非阻塞启动 HTTP
	go func() {
		if err := gs.Serve(listener); err != nil && err != http.ErrServerClosed {
			fatal(logService, fmt.Sprintf("[BOOT] HTTP server failed: %v", err.Error()))
		}
	}()
}

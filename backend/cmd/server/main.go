package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ctipscn/weauth/backend/internal/config"
	"github.com/ctipscn/weauth/backend/internal/database"
	adminHandler "github.com/ctipscn/weauth/backend/internal/handlers/admin"
	authHandler "github.com/ctipscn/weauth/backend/internal/handlers/auth"
	powHandler "github.com/ctipscn/weauth/backend/internal/handlers/pow"
	siteverifyHandler "github.com/ctipscn/weauth/backend/internal/handlers/siteverify"
	"github.com/ctipscn/weauth/backend/internal/middleware"
	"github.com/ctipscn/weauth/backend/internal/migration"
	"github.com/ctipscn/weauth/backend/internal/routers"
	"github.com/ctipscn/weauth/backend/internal/seed"
	"github.com/ctipscn/weauth/backend/internal/services/auth"
	"github.com/ctipscn/weauth/backend/internal/services/ipstrategy"
	"github.com/ctipscn/weauth/backend/internal/services/pow"
	"github.com/ctipscn/weauth/backend/internal/services/website"

	"github.com/gin-gonic/gin"
)

func main() {
	// 加载配置
	if err := config.Load(); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 连接数据库 — NewDefault 返回 *Storer 并在失败时快速失败
	dbStorePtr, err := database.NewDefault()
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}
	store := *dbStorePtr
	gormDB := store.DB()

	// 创建所有服务实例（通过构造函数注入 Storer）
	siteSvc := website.NewSiteService(store)
	ipSvc := ipstrategy.NewService(store)
	powSvc := pow.NewPowService(store, siteSvc, ipSvc)
	authSvc := auth.NewAuthService(store)

	// 创建所有控制器（依赖通过构造函数注入，不依赖全局变量）
	authCtrl := authHandler.NewAuthController(authSvc)
	powCtrl := powHandler.NewPowController(siteSvc, ipSvc, powSvc)
	siteVerifyCtrl := siteverifyHandler.NewSiteVerifyController(powSvc)
	adminCtrl := adminHandler.NewAdminController(siteSvc, ipSvc, powSvc)

	// 数据库迁移
	if err := migration.Migrate(store); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// 种子数据
	seed.SeedAdmin(gormDB)

	log.Println("Database migration completed")

	// 设置 Gin 模式
	gin.SetMode(config.AppConfig.Server.Mode)

	// 创建路由
	r := gin.Default()

	// 中间件
	r.Use(middleware.CORS())

	// 注册所有路由
	routers.RouterRegister(r, authCtrl, powCtrl, siteVerifyCtrl, adminCtrl, gormDB)

	// 创建 HTTP Server
	srv := &http.Server{
		Addr:    ":" + config.AppConfig.Server.Port,
		Handler: r,
	}

	// 在 goroutine 中启动服务器，主协程等待信号
	go func() {
		log.Printf("Server starting on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// 等待系统信号并优雅关闭
	gracefulShutdown(srv, store)
}

// gracefulShutdown 监听系统信号，优雅关闭 HTTP 服务器并释放数据库连接。
func gracefulShutdown(srv *http.Server, store database.Storer) {
	quit := make(chan os.Signal, 1)
	// 监听 SIGINT (Ctrl+C) 和 SIGTERM (systemd / Docker stop)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit

	log.Printf("Received signal %v, shutting down gracefully...", sig)

	// 设置 30 秒超时，给进行中的请求足够时间完成
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 停止接收新请求，等待进行中的请求完成
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	// 关闭数据库连接
	if err := store.Close(); err != nil {
		log.Printf("Failed to close database: %v", err)
	}

	log.Println("Server exited gracefully")
}

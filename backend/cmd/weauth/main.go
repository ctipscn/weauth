package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/ctipscn/weauth/backend/internal/bootstrap"
	"github.com/ctipscn/weauth/backend/internal/cli"
	"github.com/ctipscn/weauth/backend/internal/config"
	"github.com/ctipscn/weauth/backend/internal/infrastructure/database"
	authcredrepo "github.com/ctipscn/weauth/backend/internal/infrastructure/database/repos/auth"
	siterepo "github.com/ctipscn/weauth/backend/internal/infrastructure/database/repos/site"
	userrepo "github.com/ctipscn/weauth/backend/internal/infrastructure/database/repos/user"
	infralogger "github.com/ctipscn/weauth/backend/internal/infrastructure/logger"
)

func main() {
	// 解析命令行参数
	args, err := cli.Parse()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// 加载配置（默认值 → 配置文件 → 环境变量）
	config.Load(args.ConfigPath)

	args.Flags.ApplyFromArgs(config.Config())

	// 日志服务
	logService, err := infralogger.NewLogService(config.Config().Server.LogFile)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	logService.Info("[MAIN] Server starting...")

	// 数据库连接
	dbDriver, err := database.NewPostgresDriver(logService)
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}

	// 迁移
	bootstrap.InitDB(logService, dbDriver)

	// 用户仓储
	userRepo := userrepo.NewUserRepo(dbDriver)
	authCredRepo := authcredrepo.NewAuthCredentialRepo(dbDriver)

	// 站点仓储
	siteRepo := siterepo.NewSiteRepo(dbDriver)

	// 根据命令分流
	if args.Init {
		bootstrap.InitAdminUser(logService, userRepo, authCredRepo, args.Username, args.Password)
		return
	}

	if args.FrontendTrialURL != "" {
		bootstrap.InitTrialUser(logService, userRepo, authCredRepo, siteRepo, args.FrontendTrialURL)
		return
	}

	// 注册应用关闭
	boot := bootstrap.NewShutdown(context.Background())

	boot.RegisterCleanupFn(
		"Shutdown database",
		func(ctx context.Context) error {
			logService.Info("[MAIN] Closing database connection...")
			sqlDB, err := dbDriver.DB().DB()
			if err != nil {
				return err
			}
			return sqlDB.Close()
		},
	)

	boot.RegisterCleanupFn("Shutdown logger", func(ctx context.Context) error {
		logService.Info("[MAIN] Closing log service...")
		return logService.Close()
	})

	// 启动 HTTP 服务
	bootstrap.StartApp(logService, userRepo, authCredRepo, boot, dbDriver)

	// 等待信号
	boot.Wait()
}

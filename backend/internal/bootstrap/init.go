package bootstrap

import (
	"errors"
	"fmt"
	"log"

	"github.com/ctipscn/weauth/backend/internal/domain"
	"github.com/ctipscn/weauth/backend/internal/domain/auth"
	"github.com/ctipscn/weauth/backend/internal/domain/site"
	domainuser "github.com/ctipscn/weauth/backend/internal/domain/user"
	"github.com/ctipscn/weauth/backend/internal/infrastructure/database"
	inframodels "github.com/ctipscn/weauth/backend/internal/infrastructure/database/models"
	authcredrepo "github.com/ctipscn/weauth/backend/internal/infrastructure/database/repos/auth"
	siterepo "github.com/ctipscn/weauth/backend/internal/infrastructure/database/repos/site"
	userrepo "github.com/ctipscn/weauth/backend/internal/infrastructure/database/repos/user"
	infralogger "github.com/ctipscn/weauth/backend/internal/infrastructure/logger"
	"github.com/ctipscn/weauth/backend/internal/infrastructure/seeder"
	"github.com/ctipscn/weauth/backend/internal/infrastructure/security"

	"gorm.io/gorm"
)

// InitDB 执行数据库自动迁移并返回 *gorm.DB。
func InitDB(
	logService *infralogger.LogService,
	dbDriver *database.PostgresDriver,
) *gorm.DB {

	if err := dbDriver.AutoMigrate(
		&inframodels.User{},
		&inframodels.Auth{},
		&inframodels.Site{},
		&inframodels.Domain{},
		&inframodels.IpPolicy{},
		&inframodels.IpRecord{},
		&inframodels.Challenge{},
		&inframodels.VerificationLog{},
	); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	db := dbDriver.DB()
	logService.Info("[BOOT] Database migration completed")
	return db
}

// InitAdminUser 执行管理员用户初始化并退出。
func InitAdminUser(
	logService *infralogger.LogService,
	userRepo *userrepo.UserRepo,
	authCredRepo *authcredrepo.AuthCredentialRepo,
	username string,
	password string,
) {

	logService.Info("[MAIN] Initializing user...")

	u, err := seeder.SeedUser(userRepo, authCredRepo, seeder.SeedUserOption{
		Username: username,
		Password: password,
		Email:    "admin@example.com",
		Role:     domainuser.Administrator,
	}, logService)
	
	if err != nil {
		logService.Error(fmt.Sprintf("[MAIN] Failed to initialize user: %v", err))
		log.Fatalf("Failed to seed user: %v", err)
	}

	if u == nil {
		logService.Info("[MAIN] User already exists, skipping creation")
		fmt.Println("[MAIN] User already exists, skipping creation")
	} else {
		msg := fmt.Sprintf("[MAIN] User %s created successfully", username)
		logService.Info(msg)
		fmt.Println(msg)
	}

	logService.Info("[MAIN] User initialization completed")
}

// InitTrialUser 创建试用用户和站点后退出。
func InitTrialUser(
	logService *infralogger.LogService,
	userRepo *userrepo.UserRepo,
	authCredRepo *authcredrepo.AuthCredentialRepo,
	siteRepo *siterepo.SiteRepo,
	frontendURL string,
) {

	logService.Info("[MAIN] Initializing trial user...")

	// 检查 trial-user 是否已存在
	existingUser, err := userRepo.FindByUsername("trial-user")
	if err == nil {
		logService.Info("[MAIN] Trial user already exists, skipping creation")
		fmt.Println()
		fmt.Println("=== Trial User Already Exists ===")
		fmt.Printf("Username:     trial-user (User ID: %d)\n", existingUser.ID)
		fmt.Println()
		return
	}
	if !errors.Is(err, domain.ErrUserNotFound) {
		logService.Error(fmt.Sprintf("[MAIN] Failed to check trial user: %v", err))
		log.Fatalf("Failed to check trial user: %v", err)
	}

	// 生成随机密码
	encryptor := security.NewEncryptor()
	rawPassword, err := encryptor.GeneratePassword()
	if err != nil {
		logService.Error(fmt.Sprintf("[MAIN] Failed to generate password: %v", err))
		log.Fatalf("Failed to generate password: %v", err)
	}

	// 创建 staff 用户
	u := &domainuser.User{
		Username: "trial-user",
		Email:    "trial@example.com",
		Role:     domainuser.Staff,
		Enabled:  true,
	}
	if err := userRepo.Create(u); err != nil {
		logService.Error(fmt.Sprintf("[MAIN] Failed to create trial user: %v", err))
		log.Fatalf("Failed to create trial user: %v", err)
	}

	logService.Infof("[MAIN] Trial user created (userID=%d, username=%s)", u.ID, u.Username)

	// 哈希密码并创建认证凭证
	hash, err := encryptor.HashPassword(rawPassword)
	if err != nil {
		logService.Error(fmt.Sprintf("[MAIN] Failed to hash password: %v", err))
		log.Fatalf("Failed to hash password: %v", err)
	}

	cred := &auth.Auth{
		UserID:       u.ID,
		PasswordHash: hash,
	}
	if err := authCredRepo.Create(cred); err != nil {
		logService.Error(fmt.Sprintf("[MAIN] Failed to create auth credential: %v", err))
		log.Fatalf("Failed to create auth credential: %v", err)
	}

	// 创建站点
	s := site.NewSite(u.ID, "trial-site")
	// 覆盖默认生成的 sitekey 和 secret
	s.Sitekey = "trial-sitekey"
	s.Secret = "trial-secret"
	s.BatchModeEnabled = true

	if err := siteRepo.Create(s); err != nil {
		logService.Error(fmt.Sprintf("[MAIN] Failed to create trial site: %v", err))
		log.Fatalf("Failed to create trial site: %v", err)
	}

	// 添加域名白名单
	if _, err := siteRepo.AddDomain(s.ID, frontendURL); err != nil {
		logService.Error(fmt.Sprintf("[MAIN] Failed to add domain: %v", err))
		log.Fatalf("Failed to add domain: %v", err)
	}

	// 6. 打印结果
	fmt.Println()
	fmt.Println("=== Trial User Created ===")
	fmt.Printf("Username:     trial-user\n")
	fmt.Printf("Password:     %s\n", rawPassword)
	fmt.Printf("User ID:      %d\n", u.ID)
	fmt.Printf("Site Name:    trial-site\n")
	fmt.Printf("Sitekey:      trial-sitekey\n")
	fmt.Printf("Secret:       trial-secret\n")
	fmt.Printf("Frontend URL: %s\n", frontendURL)
	fmt.Println()

	logService.Info("[MAIN] Trial user initialization completed")
}

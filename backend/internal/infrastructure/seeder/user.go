package seeder

import (

	"github.com/ctipscn/weauth/backend/internal/domain/auth"
	domainuser "github.com/ctipscn/weauth/backend/internal/domain/user"
	infralogger "github.com/ctipscn/weauth/backend/internal/infrastructure/logger"
	"github.com/ctipscn/weauth/backend/internal/infrastructure/security"
)

// SeedUserOption 种子用户配置。
type SeedUserOption struct {
	Username string
	Password string
	Email    string
	Role     domainuser.Role
}

// SeedUser 使用领域实体创建种子用户
func SeedUser(
	userRepo domainuser.Repository,
	authCredRepo auth.Repository,
	opt SeedUserOption,
	log *infralogger.LogService,
) (*domainuser.User, error) {

	count, err := userRepo.CountAll()
	if err != nil {
		return nil, err
	}
	if count > 0 {
		log.Infof("[SEED] SeedUser: users already exist (count=%d), skipping", count)
		return nil, nil
	}

	log.Infof(
		"[SEED] SeedUser: creating user (username=%s, email=%s, role=%s)",
		opt.Username,
		opt.Email,
		opt.Role.String(),
	)

	u := &domainuser.User{
		Username: opt.Username,
		Email:    opt.Email,
		Role:     opt.Role,
		Enabled:  true,
	}
	if err := userRepo.Create(u); err != nil {
		log.Errorf("[SEED] SeedUser: failed to create user: %v", err)
		return nil, err
	}

	log.Infof("[SEED] SeedUser: hashing password")
	v := security.NewEncryptor()
	hash, err := v.HashPassword(opt.Password)
	if err != nil {
		return nil, err
	}

	cred := &auth.Auth{
		UserID:       u.ID,
		PasswordHash: hash,
	}
	if err := authCredRepo.Create(cred); err != nil {
		log.Errorf("[SEED] SeedUser: failed to create credential: %v", err)
		return nil, err
	}

	log.Infof("[SEED] SeedUser: user created (userID=%d, username=%s)", u.ID, u.Username)
	return u, nil
}


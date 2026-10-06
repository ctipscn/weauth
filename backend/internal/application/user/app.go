package user

import (
	"errors"

	"github.com/ctipscn/weauth/backend/internal/domain"
	"github.com/ctipscn/weauth/backend/internal/domain/auth"
	"github.com/ctipscn/weauth/backend/internal/domain/user"
	infralogger "github.com/ctipscn/weauth/backend/internal/infrastructure/logger"
)

// UserApp 用户应用服务
type UserApp struct {
	userRepo user.Repository
	authRepo auth.Repository
	log      *infralogger.LogService
}

// NewUserApp 创建用户应用服务
func NewUserApp(
	userRepo user.Repository, 
	authRepo auth.Repository, 
	log *infralogger.LogService,
) *UserApp {
	return &UserApp{
		userRepo: userRepo,
		authRepo: authRepo,
		log:      log,
	}
}

// GetUserByID 根据 ID 获取用户
func (app *UserApp) GetUserByID(id uint) (*user.User, error) {
	u, err := app.userRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			app.log.Warnf("[USER] GetUserByID failed: user not found (userID=%d)", id)
			return nil, ErrUserNotFound
		}
		app.log.Errorf("[USER] GetUserByID error (userID=%d): %v", id, err)
		return nil, err
	}
	return u, nil
}

// UpdateProfile 更新用户资料（用户名、头像）
func (app *UserApp) UpdateProfile(userID uint, username, avatar string) (*user.User, error) {
	u, err := app.userRepo.FindByID(userID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			app.log.Warnf("[USER] UpdateProfile failed: user not found (userID=%d)", userID)
			return nil, ErrUserNotFound
		}
		app.log.Errorf("[USER] UpdateProfile error finding user (userID=%d): %v", userID, err)
		return nil, err
	}

	oldUsername := u.Username
	if username != "" && username != u.Username {
		// 验证新用户名未被其他活跃用户使用
		existing, err := app.userRepo.FindByUsername(username)
		if err == nil && existing.ID != userID {
			app.log.Warnf("[USER] UpdateProfile failed: username already taken (userID=%d, username=%s)", userID, username)
			return nil, ErrUsernameTaken
		}
		if !errors.Is(err, domain.ErrUserNotFound) {
			app.log.Errorf("[USER] UpdateProfile error checking username (userID=%d): %v", userID, err)
			return nil, err
		}
		u.Username = username
	}
	u.Avatar = avatar

	if err := app.userRepo.Save(u); err != nil {
		app.log.Errorf("[USER] UpdateProfile DB error (userID=%d): %v", userID, err)
		return nil, err
	}

	app.log.Infof(
		"[USER] Profile updated (userID=%d, username=%s->%s, avatar=%s)", 
		userID, 
		oldUsername, 
		u.Username, 
		avatar,
	)
	return u, nil
}

// DeactivateUser 注销用户（软删除用户 + 废除认证凭证）
func (app *UserApp) DeactivateUser(userID uint) error {
	u, err := app.userRepo.FindByID(userID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			app.log.Warnf("[USER] DeactivateUser failed: user not found (userID=%d)", userID)
			return ErrUserNotFound
		}
		app.log.Errorf("[USER] DeactivateUser error finding user (userID=%d): %v", userID, err)
		return err
	}

	if !u.Enabled {
		app.log.Warnf("[USER] DeactivateUser failed: user already deactivated (userID=%d)", userID)
		return ErrUserAlreadyDeactivated
	}

	// 软删除用户（设置 enabled=false + deleted_at）
	if err := app.userRepo.Deactivate(userID); err != nil {
		app.log.Errorf("[USER] DeactivateUser DB error (userID=%d): %v", userID, err)
		return err
	}

	// 废除认证凭证（软删除 auth 记录）
	if err := app.authRepo.DeleteByUserID(userID); err != nil && !errors.Is(err, domain.ErrNotFound) {
		app.log.Errorf("[USER] DeactivateUser auth delete error (userID=%d): %v", userID, err)
		// 用户已软删除，auth 删除失败记为警告而非致命
		app.log.Warnf("[USER] DeactivateUser: auth credential not found for user (userID=%d), continuing", userID)
	}

	app.log.Infof("[USER] User deactivated (userID=%d, username=%s)", userID, u.Username)
	return nil
}

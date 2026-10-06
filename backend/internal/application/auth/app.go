package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"strings"

	"github.com/ctipscn/weauth/backend/internal/application/oidc"
	"github.com/ctipscn/weauth/backend/internal/domain"
	"github.com/ctipscn/weauth/backend/internal/domain/auth"
	"github.com/ctipscn/weauth/backend/internal/domain/user"
	infralogger "github.com/ctipscn/weauth/backend/internal/infrastructure/logger"
	"github.com/ctipscn/weauth/backend/internal/infrastructure/security"
)

// AuthApp 认证应用服务
type AuthApp struct {
	userRepo     user.Repository
	authCredRepo auth.Repository
	oidcClient   *oidc.OIDCClient
	frontendURL  string
	passwdValid  security.Encryptor
	log          *infralogger.LogService
}

// NewAuthApp 创建认证应用服务
func NewAuthApp(
	userRepo user.Repository,
	authCredRepo auth.Repository,
	oidcClient *oidc.OIDCClient,
	frontendURL string,
	passwdValid security.Encryptor,
	log *infralogger.LogService,
) *AuthApp {

	return &AuthApp{
		userRepo:     userRepo,
		authCredRepo: authCredRepo,
		oidcClient:   oidcClient,
		frontendURL:  frontendURL,
		passwdValid:  passwdValid,
		log:          log,
	}
}

// OIDCClient 暴露 OIDC 客户端
func (app *AuthApp) OIDCClient() *oidc.OIDCClient {
	return app.oidcClient
}

// FrontendURL 返回前端 URL
func (app *AuthApp) FrontendURL() string {
	return app.frontendURL
}

// FindOrCreateUserByOIDC 根据 OIDC 信息查找或创建用户
func (app *AuthApp) FindOrCreateUserByOIDC(userInfo *oidc.OIDCUserInfo) (*user.User, error) {
	// 尝试通过 OIDC Sub 查找已有凭证
	cred, err := app.authCredRepo.FindByOIDCSub(userInfo.Sub)
	if err == nil {
		// 找到已有凭证，获取关联用户
		u, err := app.userRepo.FindByID(cred.UserID)
		if err != nil {
			app.log.Errorf(
				"[AUTH] OIDC login failed to find user by credential (userID=%d): %v",
				cred.UserID,
				err,
			)
			return nil, err
		}
		if !u.Enabled {
			app.log.Warnf(
				"[AUTH] OIDC login blocked: user %d is disabled (OIDC sub=%s)",
				u.ID,
				userInfo.Sub,
			)
			return nil, ErrUserDisabled
		}

		// 不覆盖已有用户信息（用户名、头像、邮箱由用户在控制台自行维护）
		cred.UpdateLastLogin()
		if err := app.authCredRepo.Save(cred); err != nil {
			return nil, err
		}
		if err := app.userRepo.Save(u); err != nil {
			return nil, err
		}

		app.log.Infof(
			"[AUTH] OIDC login success (userID=%d, sub=%s, email=%s)",
			u.ID,
			userInfo.Sub,
			u.Email,
		)
		return u, nil
	}

	if !errors.Is(err, domain.ErrNotFound) {
		app.log.Errorf("[AUTH] OIDC lookup failed (sub=%s): %v", userInfo.Sub, err)
		return nil, err
	}

	// 创建新用户
	username := userInfo.Username
	if username == "" {
		username = userInfo.Name
	}

	u := &user.User{
		Username: username,
		Email:    userInfo.Email,
		Avatar:   userInfo.Picture,
		Role:     user.Staff,
		Enabled:  true,
	}

	if err := app.userRepo.Create(u); err != nil {
		app.log.Errorf(
			"[AUTH] OIDC user creation failed (sub=%s, email=%s): %v",
			userInfo.Sub,
			userInfo.Email,
			err,
		)
		return nil, err
	}

	// 创建认证凭证
	cred = &auth.Auth{
		UserID:    u.ID,
		OIDCSub:   &userInfo.Sub,
		LoginMode: auth.LoginModeOIDC,
	}

	if err := app.authCredRepo.Create(cred); err != nil {
		return nil, err
	}

	app.log.Infof(
		"[AUTH] OIDC new user created (userID=%d, sub=%s, email=%s)",
		u.ID,
		userInfo.Sub,
		u.Email,
	)
	return u, nil
}

// LoginByEmailOrUsername 邮箱/用户名密码登录
func (app *AuthApp) LoginByEmailOrUsername(identifier, password string) (*user.User, error) {
	u, err := app.userRepo.FindByEmail(identifier)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			u, err = app.userRepo.FindByUsername(identifier)
			if err != nil {
				if errors.Is(err, domain.ErrUserNotFound) {
					app.log.Warnf(
						"[AUTH] Login failed: user not found (identifier=%s)",
						identifier,
					)
					return nil, ErrUserNotFound
				}
				app.log.Errorf(
					"[AUTH] Login error during username lookup (identifier=%s): %v",
					identifier,
					err,
				)
				return nil, err
			}
		} else {
			app.log.Errorf(
				"[AUTH] Login error during email lookup (identifier=%s): %v",
				identifier,
				err,
			)
			return nil, err
		}
	}
	if !u.Enabled {
		app.log.Warnf(
			"[AUTH] Login blocked: user %d is disabled (email=%s)",
			u.ID,
			identifier,
		)
		return nil, ErrUserDisabled
	}

	cred, err := app.authCredRepo.FindByUsername(u.Username)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			app.log.Warnf("[AUTH] Login failed: no password set for user %s (username=%s)", u.Username, identifier)
			return nil, ErrNoPasswordSet
		}
		app.log.Errorf("[AUTH] Login error fetching credential (username=%s): %v", u.Username, err)
		return nil, err
	}

	if cred.PasswordHash == "" {
		app.log.Warnf("[AUTH] Login failed: empty password hash for user %d (email=%s)", u.ID, identifier)
		return nil, ErrNoPasswordSet
	}
	
	if !app.passwdValid.VerifyPassword(password, cred.PasswordHash) {
		app.log.Warnf("[AUTH] Login failed: invalid password for user %d (email=%s)", u.ID, identifier)
		return nil, ErrInvalidPassword
	}

	cred.UpdateLastLogin()
	if err := app.authCredRepo.Save(cred); err != nil {
		app.log.Errorf("[AUTH] Login error saving last login time (userID=%d): %v", u.ID, err)
		return nil, err
	}

	app.log.Infof("[AUTH] Login success (userID=%d, email=%s, role=%s)", u.ID, u.Email, u.Role)
	return u, nil
}

// Register 注册新用户（邮箱 + 密码，用户名默认使用邮箱）
func (app *AuthApp) Register(username, email, password string) (*user.User, error) {
	if !strings.Contains(email, "@") {
		app.log.Warnf("[AUTH] Registration failed: invalid email format (%s)", email)
		return nil, ErrInvalidEmail
	}
	if len(password) < 6 {
		app.log.Warnf("[AUTH] Registration failed: password too short (email=%s)", email)
		return nil, ErrPasswordTooShort
	}

	// 检查邮箱是否已存在
	if _, err := app.userRepo.FindByEmail(email); err == nil {
		app.log.Warnf("[AUTH] Registration failed: email already taken (email=%s)", email)
		return nil, ErrEmailTaken
	} else if !errors.Is(err, domain.ErrUserNotFound) {
		app.log.Errorf("[AUTH] Registration error checking email (email=%s): %v", email, err)
		return nil, err
	}

	// 哈希密码
	hash, err := app.passwdValid.HashPassword(password)
	if err != nil {
		app.log.Errorf("[AUTH] Registration error hashing password (email=%s): %v", email, err)
		return nil, err
	}

	// 创建用户
	u := &user.User{
		Username: username,
		Email:    email,
		Role:     user.Staff,
		Enabled:  true,
	}
	if err := app.userRepo.Create(u); err != nil {
		app.log.Errorf("[AUTH] Registration error creating user (email=%s): %v", email, err)
		return nil, err
	}

	// 创建认证凭证
	cred := &auth.Auth{
		UserID:       u.ID,
		PasswordHash: hash,
		LoginMode:    auth.LoginModePasswd,
	}
	if err := app.authCredRepo.Create(cred); err != nil {
		app.log.Errorf("[AUTH] Registration error creating credential (userID=%d): %v", u.ID, err)
		return nil, err
	}

	app.log.Infof("[AUTH] New user registered (userID=%d, email=%s)", u.ID, email)
	return u, nil
}

// UpdateUserPassword 更新用户密码（需要验证旧密码）
func (app *AuthApp) UpdateUserPassword(userID uint, oldPassword, newPassword string) error {
	cred, err := app.authCredRepo.FindByUserID(userID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			app.log.Warnf("[AUTH] Password update failed: credential not found (userID=%d)", userID)
			return ErrUserNotFound
		}
		app.log.Errorf("[AUTH] Password update error fetching credential (userID=%d): %v", userID, err)
		return err
	}

	// 有旧密码的用户必须验证旧密码
	if cred.PasswordHash != "" {
		if !app.passwdValid.VerifyPassword(oldPassword, cred.PasswordHash) {
			app.log.Warnf("[AUTH] Password update failed: old password mismatch (userID=%d)", userID)
			return ErrOldPasswordMismatch
		}
	}

	hash, err := app.passwdValid.HashPassword(newPassword)
	if err != nil {
		app.log.Errorf("[AUTH] Password update error hashing (userID=%d): %v", userID, err)
		return err
	}
	cred.PasswordHash = hash
	// OIDC-only 用户首次设置密码时，更新登录模式为 both
	if cred.LoginMode == auth.LoginModeOIDC {
		cred.LoginMode = auth.LoginModeBoth
	}
	if err := app.authCredRepo.Save(cred); err != nil {
		app.log.Errorf("[AUTH] Password update error saving (userID=%d): %v", userID, err)
		return err
	}
	app.log.Infof("[AUTH] Password updated (userID=%d)", userID)
	return nil
}



// GenerateRandomState 生成随机 state 值
func GenerateRandomState() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// SetUserPassword 为 OIDC-only 用户设置密码（无需旧密码）
func (app *AuthApp) SetUserPassword(userID uint, newPassword string) error {
	cred, err := app.authCredRepo.FindByUserID(userID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			app.log.Warnf("[AUTH] Password set failed: credential not found (userID=%d)", userID)
			return ErrUserNotFound
		}
		app.log.Errorf("[AUTH] Password set error fetching credential (userID=%d): %v", userID, err)
		return err
	}

	// 非 OIDC 登录用户不允许通过此接口设置密码
	if !cred.IsOIDCLogin() {
		app.log.Warnf("[AUTH] Password set failed: not an OIDC user (userID=%d)", userID)
		return ErrPasswordAlreadySet
	}

	hash, err := app.passwdValid.HashPassword(newPassword)
	if err != nil {
		app.log.Errorf("[AUTH] Password set error hashing (userID=%d): %v", userID, err)
		return err
	}
	cred.PasswordHash = hash
	cred.LoginMode = auth.LoginModeBoth
	if err := app.authCredRepo.Save(cred); err != nil {
		app.log.Errorf("[AUTH] Password set error saving (userID=%d): %v", userID, err)
		return err
	}
	app.log.Infof("[AUTH] Password set for OIDC user (userID=%d, loginMode=both)", userID)
	return nil
}

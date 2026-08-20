package auth

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/ctipscn/weauth/backend/internal/config"
	"github.com/ctipscn/weauth/backend/internal/database"
	"github.com/ctipscn/weauth/backend/internal/models"

	"gorm.io/gorm"
)

// AuthService 认证服务，封装 OIDC 认证与用户管理逻辑。
// 通过 Storer 接口访问数据库，便于测试和替换存储后端。
type AuthService struct {
	store            database.Storer
	cachedOIDCConfig *OIDCConfig
}

// NewAuthService 创建 AuthService 实例。
// store 必须实现 database.Storer 接口。
func NewAuthService(store database.Storer) *AuthService {
	return &AuthService{store: store}
}

// GetOIDCConfig 获取 OIDC 配置（带缓存）
func (s *AuthService) GetOIDCConfig() (*OIDCConfig, error) {
	if s.cachedOIDCConfig != nil {
		return s.cachedOIDCConfig, nil
	}

	resp, err := http.Get(config.AppConfig.OIDC.DiscoveryURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var oidcConfig OIDCConfig
	if err := json.NewDecoder(resp.Body).Decode(&oidcConfig); err != nil {
		return nil, err
	}

	s.cachedOIDCConfig = &oidcConfig
	return s.cachedOIDCConfig, nil
}

// GetOIDCAuthURL 获取 OIDC 授权 URL
func (s *AuthService) GetOIDCAuthURL(state string) (string, error) {
	oidcConfig, err := s.GetOIDCConfig()
	if err != nil {
		return "", err
	}

	params := url.Values{}
	params.Set("client_id", config.AppConfig.OIDC.ClientID)
	params.Set("redirect_uri", config.AppConfig.OIDC.RedirectURL)
	params.Set("response_type", "code")
	params.Set("scope", "openid profile email")
	params.Set("state", state)

	return oidcConfig.AuthorizationEndpoint + "?" + params.Encode(), nil
}

// ExchangeCodeForToken 用 code 换取 token
func (s *AuthService) ExchangeCodeForToken(code string) (*TokenResponse, error) {
	oidcConfig, err := s.GetOIDCConfig()
	if err != nil {
		return nil, err
	}

	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("redirect_uri", config.AppConfig.OIDC.RedirectURL)
	data.Set("client_id", config.AppConfig.OIDC.ClientID)
	data.Set("client_secret", config.AppConfig.OIDC.ClientSecret)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "POST", oidcConfig.TokenEndpoint, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, errors.New("token exchange failed: " + string(body))
	}

	var tokenResp TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, err
	}

	return &tokenResp, nil
}

// GetUserInfo 获取用户信息
func (s *AuthService) GetUserInfo(accessToken string) (*UserInfo, error) {
	oidcConfig, err := s.GetOIDCConfig()
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", oidcConfig.UserinfoEndpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, ErrOIDCUserInfoFail
	}

	var userInfo UserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, err
	}

	return &userInfo, nil
}

// FindOrCreateUserByOIDC 根据 OIDC 信息查找或创建用户
func (s *AuthService) FindOrCreateUserByOIDC(userInfo *UserInfo) (*models.User, error) {
	db := s.store.DB()
	var user models.User

	err := db.Where("oidc_sub = ?", userInfo.Sub).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 创建新用户
			user = models.User{
				OIDCSub:  userInfo.Sub,
				Username: userInfo.Username,
				Email:    userInfo.Email,
				Avatar:   userInfo.Picture,
				Role:     "user",
				Enabled:  true,
			}
			if user.Username == "" {
				user.Username = userInfo.Name
			}
			if err := db.Create(&user).Error; err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	} else {
		// 更新用户信息
		updates := map[string]any{
			"username": userInfo.Username,
			"email":    userInfo.Email,
			"avatar":   userInfo.Picture,
		}
		if userInfo.Username == "" {
			updates["username"] = userInfo.Name
		}
		db.Model(&user).Updates(updates)
	}

	if !user.Enabled {
		return nil, ErrUserDisabled
	}

	user.UpdateLastLogin()
	db.Save(&user)

	return &user, nil
}

// LoginByPassword 使用用户 ID 和密码登录
func (s *AuthService) LoginByPassword(userID uint, password string) (*models.User, error) {
	db := s.store.DB()
	var user models.User

	if err := db.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	if !user.Enabled {
		return nil, ErrUserDisabled
	}

	if !user.HasPassword() {
		return nil, ErrNoPasswordSet
	}

	if !user.CheckPassword(password) {
		return nil, ErrInvalidPassword
	}

	user.UpdateLastLogin()
	db.Save(&user)

	return &user, nil
}

// GetUserByID 根据 ID 获取用户
func (s *AuthService) GetUserByID(id uint) (*models.User, error) {
	db := s.store.DB()
	var user models.User
	if err := db.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

// UpdateUserPassword 更新用户密码
func (s *AuthService) UpdateUserPassword(userID uint, password string) error {
	db := s.store.DB()
	var user models.User
	if err := db.First(&user, userID).Error; err != nil {
		return ErrUserNotFound
	}

	if err := user.SetPassword(password); err != nil {
		return err
	}

	return db.Save(&user).Error
}

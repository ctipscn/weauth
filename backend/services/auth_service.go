package services

import (
	"github.com/ctipscn/weauth/backend/config"
	"github.com/ctipscn/weauth/backend/database"
	"github.com/ctipscn/weauth/backend/models"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"gorm.io/gorm"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserDisabled      = errors.New("user is disabled")
	ErrInvalidPassword   = errors.New("invalid password")
	ErrNoPasswordSet     = errors.New("no password set")
	ErrOIDCExchangeFail  = errors.New("failed to exchange code for token")
	ErrOIDCUserInfoFail  = errors.New("failed to get user info")
)

// OIDCConfig 缓存的 OIDC 配置
type OIDCConfig struct {
	AuthorizationEndpoint string `json:"authorization_endpoint"`
	TokenEndpoint         string `json:"token_endpoint"`
	UserinfoEndpoint      string `json:"userinfo_endpoint"`
}

var cachedOIDCConfig *OIDCConfig

// GetOIDCConfig 获取 OIDC 配置
func GetOIDCConfig() (*OIDCConfig, error) {
	if cachedOIDCConfig != nil {
		return cachedOIDCConfig, nil
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

	cachedOIDCConfig = &oidcConfig
	return cachedOIDCConfig, nil
}

// GetOIDCAuthURL 获取 OIDC 授权 URL
func GetOIDCAuthURL(state string) (string, error) {
	oidcConfig, err := GetOIDCConfig()
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

// TokenResponse OIDC token 响应
type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	IDToken     string `json:"id_token"`
}

// UserInfo OIDC 用户信息
type UserInfo struct {
	Sub      string `json:"sub"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Picture  string `json:"picture"`
}

// ExchangeCodeForToken 用 code 换取 token
func ExchangeCodeForToken(code string) (*TokenResponse, error) {
	oidcConfig, err := GetOIDCConfig()
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
func GetUserInfo(accessToken string) (*UserInfo, error) {
	oidcConfig, err := GetOIDCConfig()
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
func FindOrCreateUserByOIDC(userInfo *UserInfo) (*models.AdminUser, error) {
	var user models.AdminUser

	err := database.DB.Where("oidc_sub = ?", userInfo.Sub).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 创建新用户
			user = models.AdminUser{
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
			if err := database.DB.Create(&user).Error; err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	} else {
		// 更新用户信息
		updates := map[string]interface{}{
			"username": userInfo.Username,
			"email":    userInfo.Email,
			"avatar":   userInfo.Picture,
		}
		if userInfo.Username == "" {
			updates["username"] = userInfo.Name
		}
		database.DB.Model(&user).Updates(updates)
	}

	if !user.Enabled {
		return nil, ErrUserDisabled
	}

	user.UpdateLastLogin()
	database.DB.Save(&user)

	return &user, nil
}

// LoginByPassword 使用用户 ID 和密码登录
func LoginByPassword(userID uint, password string) (*models.AdminUser, error) {
	var user models.AdminUser

	if err := database.DB.First(&user, userID).Error; err != nil {
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
	database.DB.Save(&user)

	return &user, nil
}

// GetUserByID 根据 ID 获取用户
func GetUserByID(id uint) (*models.AdminUser, error) {
	var user models.AdminUser
	if err := database.DB.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

// UpdateUserPassword 更新用户密码
func UpdateUserPassword(userID uint, password string) error {
	var user models.AdminUser
	if err := database.DB.First(&user, userID).Error; err != nil {
		return ErrUserNotFound
	}

	if err := user.SetPassword(password); err != nil {
		return err
	}

	return database.DB.Save(&user).Error
}

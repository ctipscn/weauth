package oidc

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

// OIDCClient OIDC HTTP 客户端
type OIDCClient struct {
	mu           sync.RWMutex
	discoveryURL string
	cachedConfig *OIDCConfig
	httpClient   *http.Client
}

// NewOIDCClient 创建 OIDC HTTP 客户端
func NewOIDCClient() *OIDCClient {
	return &OIDCClient{
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

// SetDiscoveryURL 设置 OIDC 发现 URL
func (c *OIDCClient) SetDiscoveryURL(url string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.discoveryURL = url
	c.cachedConfig = nil
}

// GetConfig 获取 OIDC 配置（带缓存）
func (c *OIDCClient) GetConfig() (*OIDCConfig, error) {
	c.mu.RLock()
	if c.cachedConfig != nil {
		defer c.mu.RUnlock()
		return c.cachedConfig, nil
	}
	c.mu.RUnlock()

	c.mu.Lock()
	defer c.mu.Unlock()

	// 双重检查
	if c.cachedConfig != nil {
		return c.cachedConfig, nil
	}

	if c.discoveryURL == "" {
		return nil, errors.New("OIDC discovery URL not set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", c.discoveryURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var oidcConfig OIDCConfig
	if err := json.NewDecoder(resp.Body).Decode(&oidcConfig); err != nil {
		return nil, err
	}

	c.cachedConfig = &oidcConfig
	return c.cachedConfig, nil
}

// GetAuthURL 获取 OIDC 授权 URL
func (c *OIDCClient) GetAuthURL(clientID, redirectURL, state string) (string, error) {
	oidcConfig, err := c.GetConfig()
	if err != nil {
		return "", err
	}

	params := url.Values{}
	params.Set("client_id", clientID)
	params.Set("redirect_uri", redirectURL)
	params.Set("response_type", "code")
	params.Set("scope", "openid profile email")
	params.Set("state", state)

	return oidcConfig.AuthorizationEndpoint + "?" + params.Encode(), nil
}

// ExchangeCode 用 code 换取 token
func (c *OIDCClient) ExchangeCode(clientID, clientSecret, redirectURL, code string) (*TokenResponse, error) {
	oidcConfig, err := c.GetConfig()
	if err != nil {
		return nil, err
	}

	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("redirect_uri", redirectURL)
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "POST", oidcConfig.TokenEndpoint, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.httpClient.Do(req)
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
func (c *OIDCClient) GetUserInfo(accessToken string) (*OIDCUserInfo, error) {
	oidcConfig, err := c.GetConfig()
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

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, ErrOIDCUserInfoFail
	}

	var userInfo OIDCUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, err
	}

	return &userInfo, nil
}
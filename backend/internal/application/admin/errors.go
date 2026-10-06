// Package admin 提供管理员后台应用服务
package admin

import "errors"

var (
	// ErrUserNotFound 用户不存在
	ErrUserNotFound = errors.New("user not found")
	// ErrSiteNotFound 站点不存在
	ErrSiteNotFound = errors.New("site not found")
	// ErrDomainNotFound 域名不存在
	ErrDomainNotFound = errors.New("domain not found")
)
package website

import "errors"

var (
	ErrSiteNotFound     = errors.New("site not found")
	ErrSiteDisabled     = errors.New("site is disabled")
	ErrInvalidSitekey   = errors.New("invalid sitekey")
	ErrInvalidSecret    = errors.New("invalid secret")
	ErrDomainNotFound   = errors.New("domain not found")
	ErrDomainMismatch   = errors.New("domain not in whitelist")
	ErrPermissionDenied = errors.New("permission denied")
)
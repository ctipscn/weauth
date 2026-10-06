package domain

import "errors"

var (
	// 通用错误
	ErrNotFound     = errors.New("not found")
	ErrUnauthorized = errors.New("unauthorized")
	ErrForbidden    = errors.New("forbidden")

	// 用户错误
	ErrUserNotFound    = errors.New("user not found")
	ErrUserDisabled    = errors.New("user is disabled")
	ErrInvalidPassword = errors.New("invalid password")
	ErrNoPasswordSet   = errors.New("no password set")

	// 站点错误
	ErrSiteNotFound   = errors.New("site not found")
	ErrSiteDisabled   = errors.New("site is disabled")
	ErrInvalidSitekey = errors.New("invalid sitekey")
	ErrInvalidSecret  = errors.New("invalid secret")

	// 域名错误
	ErrDomainNotFound      = errors.New("domain not found")
	ErrDomainMismatch      = errors.New("domain not in whitelist")
	ErrDomainAlreadyExists = errors.New("domain already exists")

	// 挑战错误
	ErrChallengeNotFound = errors.New("challenge not found")
	ErrChallengeExpired  = errors.New("challenge expired")
	ErrChallengeSolved   = errors.New("challenge already solved")
	ErrInvalidSolution   = errors.New("invalid solution")
	ErrTokenNotFound     = errors.New("token not found")
	ErrTokenUsed         = errors.New("token already used")
	ErrTokenExpired      = errors.New("token expired")
	ErrNoncesMismatch    = errors.New("nonces count mismatch")

	// OIDC 错误
	ErrOIDCExchangeFail = errors.New("failed to exchange code for token")
	ErrOIDCUserInfoFail = errors.New("failed to get user info")
)
package pow

import "errors"

// PoW 相关错误
var (
	ErrChallengeNotFound = errors.New("challenge not found")
	ErrChallengeExpired  = errors.New("challenge expired")
	ErrChallengeSolved   = errors.New("challenge already solved")
	ErrInvalidSolution   = errors.New("invalid solution")
	ErrTokenNotFound     = errors.New("token not found")
	ErrTokenUsed         = errors.New("token already used")
	ErrTokenExpired      = errors.New("token expired")
	ErrNoncesMismatch    = errors.New("nonces count mismatch")
)

// 站点相关错误（由 pow 包统一提供，避免 site ↔ pow 循环依赖）
var (
	ErrSiteNotFound     = errors.New("site not found")
	ErrSiteDisabled     = errors.New("site is disabled")
	ErrInvalidSitekey   = errors.New("invalid sitekey")
	ErrInvalidSecret    = errors.New("invalid secret")
	ErrDomainNotFound   = errors.New("domain not found")
	ErrDomainMismatch   = errors.New("domain not in whitelist")
	ErrPermissionDenied = errors.New("permission denied")
)
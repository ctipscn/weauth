package auth

import "time"

type Auth struct {
	ID           uint
	UserID       uint
	OIDCSub      *string
	PasswordHash string
	LastLoginAt  *time.Time
	LoginMode    LoginMode
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (a *Auth) IsOIDCLogin() bool {
	return a.OIDCSub != nil && *a.OIDCSub != ""
}

func (a *Auth) UpdateLastLogin() {
	now := time.Now()
	a.LastLoginAt = &now
}

package auth

import "errors"

var (
	ErrUserNotFound     = errors.New("user not found")
	ErrUserDisabled     = errors.New("user is disabled")
	ErrInvalidPassword  = errors.New("invalid password")
	ErrNoPasswordSet    = errors.New("no password set")
	ErrOIDCExchangeFail = errors.New("failed to exchange code for token")
	ErrOIDCUserInfoFail = errors.New("failed to get user info")
)
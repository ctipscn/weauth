package oidc

import "errors"

var (
	ErrOIDCExchangeFail = errors.New("failed to exchange code for token")
	ErrOIDCUserInfoFail = errors.New("failed to get user info")
)
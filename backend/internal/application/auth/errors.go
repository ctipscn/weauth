package auth

import "errors"

var (
	ErrUserNotFound         = errors.New("user not found")
	ErrUserDisabled         = errors.New("user is disabled")
	ErrInvalidPassword      = errors.New("invalid password")
	ErrNoPasswordSet        = errors.New("no password set")
	ErrEmailTaken           = errors.New("email already taken")
	ErrInvalidEmail          = errors.New("invalid email")
	ErrPasswordTooShort      = errors.New("password too short")
	ErrOldPasswordMismatch   = errors.New("old password is incorrect")
	ErrPasswordAlreadySet    = errors.New("password already set")
)
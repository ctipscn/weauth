package user

import "errors"

var (
	ErrUserNotFound          = errors.New("user not found")
	ErrUsernameTaken         = errors.New("username already taken")
	ErrEmailTaken            = errors.New("email already taken")
	ErrUserAlreadyDeactivated = errors.New("user already deactivated")
)
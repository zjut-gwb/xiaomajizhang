package errors

import "errors"

var (
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrUserDisabled        = errors.New("user disabled")
	ErrUserNotFound        = errors.New("user not found")
	ErrUsernameExists      = errors.New("username exists")
	ErrEmailExists         = errors.New("email exists")
	ErrRefreshTokenInvalid = errors.New("refresh token invalid")
)

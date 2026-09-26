package user

import "errors"

var (
	ErrInvalidCredentials = errors.New("invalid phone or password")
	ErrTokenRevoked       = errors.New("refresh token is revoked or expired")
)

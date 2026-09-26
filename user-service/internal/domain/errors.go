package domain

import "errors"

var (
	ErrNotFound     = errors.New("user not found")
	ErrConflict     = errors.New("user with this phone already exists")
	ErrInvalidInput = errors.New("invalid user data")
)

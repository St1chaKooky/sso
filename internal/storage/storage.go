package storage

import "errors"

var (
	ErrAppNotFound  = errors.New("app not found")
	ErrUserNotFound = errors.New("user not found")
	ErrUserExists   = errors.New("user already exists")
)

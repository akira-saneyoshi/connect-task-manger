package model

import "errors"

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrAuthentication    = errors.New("authentication failed")
	ErrInvalidToken      = errors.New("invalid token")
	ErrUnauthorized      = errors.New("unauthorized")
)

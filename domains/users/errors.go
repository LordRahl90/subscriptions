package users

import "errors"

var (
	// ErrUserExists user exists error
	ErrUserExists = errors.New("email already exists")
)

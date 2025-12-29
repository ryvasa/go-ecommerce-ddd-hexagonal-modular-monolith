package error

import "errors"

var (
	ErrEmailAlreadyUsed  = errors.New("email already used")
	ErrInvalidCredential = errors.New("invalid credentials")
	ErrUserNotFound      = errors.New("user not found")
)

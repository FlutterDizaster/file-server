package apperrors

import (
	"errors"
)

var (
	// User management errors.
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrWrongCredentials  = errors.New("wrong credentials")

	// Documents management errors.
	ErrWrongMetadata      = errors.New("wrong metadata")
	ErrWrongFilterOptions = errors.New("wrong filter options")
	ErrAccessDenied       = errors.New("access denied")
	ErrNotFound           = errors.New("not found")
)

type Error struct {
	Code    int
	Message string
}

func (e Error) Error() string {
	return e.Message
}

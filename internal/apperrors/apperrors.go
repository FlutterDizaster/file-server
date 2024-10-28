package apperrors

import (
	"net/http"
)

var (
	// User management errors.
	ErrUserAlreadyExists = Error{
		Code:    http.StatusConflict,
		Message: "user already exists",
	}
	ErrWrongCredentials = Error{
		Code:    http.StatusUnauthorized,
		Message: "wrong credentials",
	}

	// General.
	ErrAccessDenied = Error{
		Code:    http.StatusForbidden,
		Message: "access denied",
	}
	ErrNotFound = Error{
		Code:    http.StatusNotFound,
		Message: "not found",
	}

	// Documents management errors.
	ErrWrongMetadata = Error{
		Code:    http.StatusBadRequest,
		Message: "wrong metadata",
	}
	ErrWrongFilterOptions = Error{
		Code:    http.StatusBadRequest,
		Message: "wrong filter options",
	}

	// HTTP errors.
	ErrInvalidContentType = Error{
		Code:    http.StatusBadRequest,
		Message: "invalid content type",
	}
	ErrAuthorizationHeaderNotFound = Error{
		Code:    http.StatusUnauthorized,
		Message: "authorization header not found",
	}
	ErrInvalidToken = Error{
		Code:    http.StatusUnauthorized,
		Message: "invalid token",
	}
)

type Error struct {
	Code    int
	Message string
}

func (e Error) Error() string {
	return e.Message
}

func NewBadCredentialError(message string) Error {
	return Error{
		Code:    http.StatusUnauthorized,
		Message: message,
	}
}

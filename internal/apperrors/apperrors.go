package apperrors

import (
	"net/http"
)

var (
	// User management errors.

	// User alredy exists.
	ErrUserAlreadyExists = Error{
		Code:    http.StatusConflict,
		Message: "user already exists",
	}
	// Wrong credentials.
	ErrWrongCredentials = Error{
		Code:    http.StatusUnauthorized,
		Message: "wrong credentials",
	}

	// General.

	// Access denied.
	ErrAccessDenied = Error{
		Code:    http.StatusForbidden,
		Message: "access denied",
	}
	// Not found.
	ErrNotFound = Error{
		Code:    http.StatusNotFound,
		Message: "not found",
	}

	// Documents management errors.

	// Wrong metadata.
	ErrWrongMetadata = Error{
		Code:    http.StatusBadRequest,
		Message: "wrong metadata",
	}
	// Wrong filter options.
	ErrWrongFilterOptions = Error{
		Code:    http.StatusBadRequest,
		Message: "wrong filter options",
	}
	// Unknown filter name.
	ErrUnknownFilter = Error{
		Code:    http.StatusBadRequest,
		Message: "unknown filter name",
	}
	// Invalid filter value.
	ErrInvalidFilterValue = Error{
		Code:    http.StatusBadRequest,
		Message: "invalid filter value",
	}

	// HTTP errors.

	// Invalid content type.
	ErrInvalidContentType = Error{
		Code:    http.StatusBadRequest,
		Message: "invalid content type",
	}
	// Invalid request body.
	ErrAuthorizationHeaderNotFound = Error{
		Code:    http.StatusUnauthorized,
		Message: "authorization header not found",
	}
	// Invalid token.
	ErrInvalidToken = Error{
		Code:    http.StatusUnauthorized,
		Message: "invalid token",
	}
)

// Error is a custom error type.
// Code is an HTTP status code.
// Message is an error message.
type Error struct {
	Code    int
	Message string
}

// Error implements the error interface.
func (e Error) Error() string {
	return e.Message
}

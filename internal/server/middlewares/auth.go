package middlewares

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/FlutterDizaster/file-server/internal/apperrors"
	jwtresolver "github.com/FlutterDizaster/file-server/internal/jwt-resolver"
	"github.com/FlutterDizaster/file-server/internal/models"
)

type CtxKey int

const (
	KeyUserID CtxKey = iota
)

// Auth is a stateful middleware that checks if user is authorized.
// If user is authorized, it adds user ID to the requests context.
// Otherwise, it returns an error.
type Auth struct {
	Resolver *jwtresolver.JWTResolver
}

// Handle method handles incoming requests.
func (a *Auth) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Try to get token from Authorization header
		token := r.Header.Get("Authorization")

		// If token not found, return error
		// Otherwise, try to decode it
		if token == "" {
			a.responseWithError(w, r, apperrors.ErrAuthorizationHeaderNotFound)
			return
		}

		// Try to decode token
		claims, err := a.Resolver.DecryptToken(token)
		if err != nil {
			a.responseWithError(w, r, apperrors.ErrInvalidToken)
			return
		}

		// Add user ID to context
		ctx := context.WithValue(r.Context(), KeyUserID, claims.UserID)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func (a Auth) responseWithError(w http.ResponseWriter, r *http.Request, err error) {
	resp := &models.Response{
		Error: &models.ResponseError{},
	}
	var appserror *apperrors.Error

	switch {
	case errors.As(err, &appserror):
		resp.Error.Code = appserror.Code
		resp.Error.Text = appserror.Message
	default:
		slog.Error(
			"Error while processing request",
			slog.String("Method", r.Method),
			slog.String("URL", r.URL.String()),
			slog.Any("err", err),
		)
		resp.Error.Code = http.StatusInternalServerError
		resp.Error.Text = err.Error()
	}

	respData, err := resp.MarshalJSON()
	if err != nil {
		slog.Error("Error while marshaling response", slog.Any("err", err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.Error.Code)

	if _, err = w.Write(respData); err != nil {
		slog.Error("Error while writing response", slog.Any("err", err))
		return
	}
}

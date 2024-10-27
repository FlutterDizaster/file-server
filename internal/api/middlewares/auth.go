package middlewares

import (
	"context"
	"log/slog"
	"net/http"

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
		var respErr *models.ResponseError

		// Try to get token from Authorization header
		token := r.Header.Get("Authorization")

		// If token not found, return error
		// Otherwise, try to decode it
		if token == "" {
			respErr = &models.ResponseError{
				Code: http.StatusUnauthorized,
				Text: "Authorization header not found",
			}
		} else {
			// Try to decode token
			claims, err := a.Resolver.DecryptToken(token)
			if err != nil {
				respErr = &models.ResponseError{
					Code: http.StatusUnauthorized,
					Text: err.Error(),
				}
			} else {
				// Add user ID to context
				ctx := context.WithValue(r.Context(), KeyUserID, claims.UserID)
				r = r.WithContext(ctx)
			}
		}

		// If error occured, return it
		// Otherwise, return next handler
		if respErr != nil {
			slog.Info("Failed to authorize request", slog.Any("err", respErr.Text))

			w.WriteHeader(respErr.Code)
			w.Header().Set("Content-Type", "application/json")

			// Create response
			resp := &models.Response{
				Error: respErr,
			}

			// Marshal response
			respData, err := resp.MarshalJSON()
			if err != nil {
				slog.Error("Error while creating response", slog.Any("err", err))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			// Write response
			if _, err = w.Write(respData); err != nil {
				slog.Error("Error while writing response", slog.Any("err", err))
			}
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

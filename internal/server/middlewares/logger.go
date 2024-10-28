package middlewares

import (
	"log/slog"
	"net/http"
	"time"
)

// loggerWriter is a wrapper for http.ResponseWriter that keeps track of the status code.
type loggerWriter struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader stores the status code and calls the underlying ResponseWriter's WriteHeader method.
func (w *loggerWriter) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

// Logger is a middleware that logs incoming requests
// and logging some of the details to slog.
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		lw := loggerWriter{w, http.StatusOK}
		next.ServeHTTP(&lw, r)

		slog.Info(
			"Incoming request",
			slog.Int("status", lw.statusCode),
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.Duration("duration", time.Since(start)),
			slog.String("client ip", getClientIP(r)),
		)
	})
}

// getClientIP helper for getting client IP.
func getClientIP(r *http.Request) string {
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		return ip
	}
	if ip := r.Header.Get("X-Real-Ip"); ip != "" {
		return ip
	}
	return r.RemoteAddr
}

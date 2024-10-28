package middlewares

import "net/http"

// Middleware is a function that takes an http.Handler and returns an http.Handler.
type Middleware func(http.Handler) http.Handler

// MakeChain wraps a list of middlewares in a chain.
// Chain works like FIFO stack of middlewares.
func MakeChain(mws ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for _, mw := range mws {
			next = mw(next)
		}
		return next
	}
}

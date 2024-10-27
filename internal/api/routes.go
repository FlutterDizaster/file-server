package api

import (
	"net/http"

	"github.com/FlutterDizaster/file-server/internal/api/middlewares"
)

// setupRouter sets up the router with the appropriate middleware chains
// and returns the router.
func (a *API) setupRouter() *http.ServeMux {
	router := http.NewServeMux()

	// Public routes
	userRouter := http.NewServeMux()
	userRouter.HandleFunc("POST /auth", a.userAuthHandler)
	userRouter.HandleFunc("POST /register", a.userRegisterHandler)

	// Private routes
	docRouter := http.NewServeMux()
	docRouter.HandleFunc("GET /{id}", a.docGetHandler)
	docRouter.HandleFunc("HEAD /{id}", a.docGetHeadHandler)
	docRouter.HandleFunc("GET /", a.docGetListHandler)
	docRouter.HandleFunc("HEAD /", a.docGetListHeadHandler)
	docRouter.HandleFunc("POST /", a.docPostHandler)
	docRouter.HandleFunc("DELETE /{id}", a.docDeleteHandler)

	// Public middleware chain
	publicChain := middlewares.MakeChain(
		middlewares.Logger,
	)

	// Private middleware chain
	authMw := middlewares.Auth{
		Resolver: a.jwtResolver,
	}
	privateChain := middlewares.MakeChain(
		middlewares.Logger,
		authMw.Handle,
	)

	// Setup general router
	router.Handle("/api/", publicChain(http.StripPrefix("/api/", userRouter)))
	router.Handle("/api/docs", privateChain(http.StripPrefix("/api/docs", docRouter)))

	return router
}

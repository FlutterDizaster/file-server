package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"golang.org/x/sync/errgroup"
)

// Settings used to create Server.
// Settings must be provided to New function.
// All fields are required and cant be nil.
type Settings struct {
	Addr string
	Port string

	Handler http.Handler

	ShutdownMaxTime time.Duration
}

// Server is a http server.
// Can be started with Start method.
// Gracefully stops when context is canceled.
// Must be initialized with New function.
type Server struct {
	server          *http.Server
	addr            string
	port            string
	shutdownMaxTime time.Duration
}

// New creates new Server instance.
// It takes Settings as argument and returns pointer to Server.
// All fields of Settings are required and cant be nil.
func New(settings Settings) *Server {
	a := &Server{
		addr: settings.Addr,
		port: settings.Port,

		shutdownMaxTime: settings.ShutdownMaxTime,
	}

	a.server = &http.Server{
		ReadHeaderTimeout: time.Second,
		Addr:              fmt.Sprintf("%s:%s", a.addr, a.port),
		Handler:           settings.Handler,
	}

	return a
}

// Start starts the server and blocks until context is canceled.
// It returns error if server cant be started.
// Server is gracefully stopped when context is canceled.
// If error occurs on start, shutdown is skipped.
func (a *Server) Start(ctx context.Context) error {
	eg, egCtx := errgroup.WithContext(ctx)

	errorOnStart := false

	eg.Go(func() error {
		<-egCtx.Done()

		// Skip shutdown if error on start
		if errorOnStart {
			return nil
		}

		slog.Info("Shutting down server")
		shutdownCtx, cancle := context.WithTimeout(context.Background(), a.shutdownMaxTime)
		defer cancle()

		return a.server.Shutdown(shutdownCtx)
	})

	eg.Go(func() error {
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Error while starting server", slog.Any("err", err))
			errorOnStart = true
			return err
		}
		return nil
	})

	return eg.Wait()
}

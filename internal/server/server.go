package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"golang.org/x/sync/errgroup"
)

type Settings struct {
	Addr string
	Port string

	Handler http.Handler

	ShutdownMaxTime time.Duration
}

type Server struct {
	server          *http.Server
	addr            string
	port            string
	shutdownMaxTime time.Duration
}

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

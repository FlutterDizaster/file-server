package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/FlutterDizaster/file-server/internal/application"
)

func main() {
	os.Exit(mainWithCode())
}

func mainWithCode() int {
	// Gracefull shutdown with SIGINT and SIGTERM
	ctx, cancle := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	defer cancle()

	app, err := application.New(ctx)
	if err != nil {
		slog.Error("Error while creating application", slog.Any("err", err))
		return 1
	}

	err = app.Start(ctx)
	if err != nil {
		slog.Error("Error while starting application", slog.Any("err", err))
		return 1
	}

	return 0
}

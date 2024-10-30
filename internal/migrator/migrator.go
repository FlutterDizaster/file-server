package migrator

import (
	"context"
	"errors"
	"log/slog"

	"github.com/golang-migrate/migrate/v4"

	//nolint:revieve // This is for migrate
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5"
)

func RunMigrations(ctx context.Context, connStr, migrationsPath string) error {
	if migrationsPath == "" {
		slog.Info("Skipping migrations")
		return nil
	}

	slog.Info("Running migrations")

	migrator, err := migrate.New(
		"file://"+migrationsPath,
		connStr,
	)
	if err != nil {
		slog.Error("Error while creating migrator", slog.Any("err", err))
		return err
	}
	defer migrator.Close()

	errCh := make(chan error)

	go func() {
		errCh <- migrator.Up()
	}()

	select {
	case err = <-errCh:
		if err != nil && !errors.Is(err, migrate.ErrNoChange) {
			slog.Error("Error while running migrations", slog.Any("err", err))
			return err
		}
	case <-ctx.Done():
		migrator.GracefulStop <- true
		return ctx.Err()
	}

	slog.Info("Migrations finished")

	return nil
}

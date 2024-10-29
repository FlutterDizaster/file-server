package postgresrepo

import (
	"context"
	"errors"
	"log/slog"

	"github.com/golang-migrate/migrate/v4"
	"github.com/jackc/pgx/v5/pgxpool"

	//nolint:revieve // This is for migrate
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5"
)

// TODO: Add retry functionality.

type Settings struct {
	ConnectionString string
	MigrationsPath   string
}

type PostgresRepository struct {
	pool           *pgxpool.Pool
	config         *pgxpool.Config
	connStr        string
	migrationsPath string
}

func New(ctx context.Context, settings Settings) (*PostgresRepository, error) {
	repo := &PostgresRepository{
		connStr:        settings.ConnectionString,
		migrationsPath: settings.MigrationsPath,
	}

	// Parse connection string
	config, err := pgxpool.ParseConfig(repo.connStr)
	if err != nil {
		slog.Error("Error while parsing connection string", slog.Any("err", err))
		return nil, err
	}
	repo.config = config

	// Run migrations
	err = repo.runMigrations(ctx)
	if err != nil {
		slog.Error("Error while running migrations", slog.Any("err", err))
		return nil, err
	}

	// Connect to database
	err = repo.connect(ctx)
	if err != nil {
		slog.Error("Error while connecting to database", slog.Any("err", err))
		return nil, err
	}

	return repo, nil
}

func (p *PostgresRepository) connect(ctx context.Context) error {
	pool, err := pgxpool.NewWithConfig(ctx, p.config)
	if err != nil {
		return err
	}

	p.pool = pool

	return nil
}

func (p *PostgresRepository) runMigrations(ctx context.Context) error {
	if p.migrationsPath == "" {
		slog.Info("Skipping migrations")
		return nil
	}

	slog.Info("Running migrations")

	migrator, err := migrate.New(
		"file://"+p.migrationsPath,
		p.connStr,
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

package postgresrepo

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

// PostgresRepository is a repository for PostgreSQL database.
// It used to upload, download and delete metadata.
type PostgresRepository struct {
	pool    *pgxpool.Pool
	config  *pgxpool.Config
	connStr string
}

// New creates a new instance of PostgresRepository using the provided
// connection string. It parses the connection string to configure
// the database connection and then establishes the connection.
// Returns a pointer to the created PostgresRepository and an error if
// any step fails.
func New(ctx context.Context, connStr string) (*PostgresRepository, error) {
	repo := &PostgresRepository{
		connStr: connStr,
	}

	// Parse connection string
	config, err := pgxpool.ParseConfig(repo.connStr)
	if err != nil {
		slog.Error("Error while parsing connection string", slog.Any("err", err))
		return nil, err
	}
	repo.config = config

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

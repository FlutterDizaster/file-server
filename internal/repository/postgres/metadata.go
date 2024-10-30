package postgresrepo

import (
	"context"
	"errors"
	"log/slog"
	"strings"
	"time"

	"github.com/FlutterDizaster/file-server/internal/models"
	"github.com/google/uuid"
)

func (p PostgresRepository) UploadMetadata(ctx context.Context, meta models.Metadata) error {
	// Start transaction
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		slog.Error("Error while starting transaction", slog.Any("err", err))
		return err
	}

	// Add metadata to metadata table
	row := tx.QueryRow(
		ctx,
		queryUploadMetadata,
		meta.Name,
		meta.File,
		meta.Public,
		meta.Mime,
		meta.OwnerID,
		meta.JSON,
		meta.FileSize,
	)

	var id uuid.UUID
	err = row.Scan(&id)

	if err != nil {
		slog.Error("Error while inserting metadata", slog.Any("err", err))
		return errors.Join(err, tx.Rollback(ctx))
	}

	// Add users to meta_access table
	for _, login := range meta.Grant {
		_, err = tx.Exec(ctx, queryGrantMetadataAcsess, id, login)
		if err != nil {
			slog.Error("Error while inserting access grant", slog.Any("err", err))
			return errors.Join(err, tx.Rollback(ctx))
		}
	}

	// Commit transaction
	err = tx.Commit(ctx)
	if err != nil {
		slog.Error("Error while committing transaction", slog.Any("err", err))
		return errors.Join(err, tx.Rollback(ctx))
	}

	return nil
}

func (p PostgresRepository) GetMetadataByUserID(
	ctx context.Context,
	userID uuid.UUID,
) ([]models.Metadata, error) {
	rows, err := p.pool.Query(
		ctx,
		queryGetUsersMetadata,
		userID,
	)
	if err != nil {
		return nil, err
	}

	var metaList []models.Metadata

	for rows.Next() {
		var (
			meta        models.Metadata
			grantStr    string
			createdTime time.Time
		)

		err = rows.Scan(
			&meta.ID,
			&meta.Name,
			&meta.Mime,
			&meta.File,
			&meta.Public,
			&createdTime,
			&meta.OwnerID,
			&meta.JSON,
			&meta.FileSize,
			&grantStr,
		)
		if err != nil {
			return nil, err
		}

		meta.Created = createdTime.Format(time.DateTime)

		meta.Grant = strings.Split(grantStr, ",")

		metaList = append(metaList, meta)
	}

	return metaList, nil
}

func (p PostgresRepository) DeleteMetadata(ctx context.Context, id, userID uuid.UUID) error {
	_, err := p.pool.Exec(ctx, queryDeleteMetadata, id, userID)
	return err
}

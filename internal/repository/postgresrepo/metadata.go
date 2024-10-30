package postgresrepo

import (
	"context"
	"log/slog"
	"strings"
	"time"

	"github.com/FlutterDizaster/file-server/internal/models"
	"github.com/google/uuid"
)

// UploadMetadata uploads metadata to the PostgreSQL database.
//
// It begins a transaction, inserts metadata into the metadata table,
// and grants access to specified users by adding entries to the meta_access table.
//
// If any step fails, the transaction is rolled back, and an error is returned.
//
// Returns the UUID of the newly inserted metadata if successful, or an error if not.
func (p PostgresRepository) UploadMetadata(
	ctx context.Context,
	meta models.Metadata,
) (uuid.UUID, error) {
	// Start transaction
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		slog.Error("Error while starting transaction", slog.Any("err", err))
		return uuid.Nil, err
	}

	defer func(e error) {
		if e != nil {
			//nolint:errcheck // ignore
			tx.Rollback(ctx)
		}
	}(err)

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
		return uuid.Nil, err
	}

	// Add users to meta_access table
	for _, login := range meta.Grant {
		_, err = tx.Exec(ctx, queryGrantMetadataAcsess, id, login)
		if err != nil {
			slog.Error("Error while inserting access grant", slog.Any("err", err))
			return uuid.Nil, err
		}
	}

	// Commit transaction
	err = tx.Commit(ctx)
	if err != nil {
		slog.Error("Error while committing transaction", slog.Any("err", err))
		return uuid.Nil, err
	}

	return id, nil
}

// GetMetadataByUserID retrieves metadata associated with a given user ID from the PostgreSQL database.
//
// It queries the metadata table to fetch all metadata records belonging to the specified user ID.
// Each record includes information such as ID, name, MIME type, file status, public visibility,
// creation time, owner ID, JSON data, file size, and access grants.
//
// Returns a slice of models.Metadata if successful, or an error if the query fails or if there is an issue
// scanning the rows.
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

	defer rows.Close()

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return metaList, nil
}

// DeleteMetadata delete metadata from repository.
// Returns error if delete failed.
func (p PostgresRepository) DeleteMetadata(ctx context.Context, id, userID uuid.UUID) error {
	_, err := p.pool.Exec(ctx, queryDeleteMetadata, id, userID)
	return err
}

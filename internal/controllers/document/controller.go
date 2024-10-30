package docctrl

import (
	"context"
	"errors"
	"io"

	"github.com/FlutterDizaster/file-server/internal/apperrors"
	"github.com/FlutterDizaster/file-server/internal/docfilter"
	"github.com/FlutterDizaster/file-server/internal/models"
	"github.com/google/uuid"
)

type FileRepository interface {
	UploadFile(ctx context.Context, file io.Reader, fileSize int64) error
	GetFile(ctx context.Context, meta models.Metadata) (io.ReadSeekCloser, error)
	DeleteFile(ctx context.Context, id string) error
}

type MetadataRepository interface {
	UploadMetadata(ctx context.Context, meta models.Metadata) error
	GetMetadataByUserID(ctx context.Context, userID uuid.UUID) ([]models.Metadata, error)
	DeleteMetadata(ctx context.Context, id, userID uuid.UUID) error
}

type UserRepository interface {
	GetUserByLogin(ctx context.Context, login string) (models.User, error)
}

type MetadataCache interface {
	InvalidateUserCache(ctx context.Context, id uuid.UUID) error
	SaveUserCache(ctx context.Context, id uuid.UUID, meta []models.Metadata) error
	GetUserCache(ctx context.Context, id uuid.UUID) ([]models.Metadata, error)
}

type Settings struct {
	FileRepo FileRepository
	MetaRepo MetadataRepository
	UserRepo UserRepository
	Cache    MetadataCache
}

type DocumentsController struct {
	fileRepo FileRepository
	metaRepo MetadataRepository
	userRepo UserRepository
	cache    MetadataCache
}

func New(settings Settings) *DocumentsController {
	ctrl := &DocumentsController{
		fileRepo: settings.FileRepo,
		metaRepo: settings.MetaRepo,
		userRepo: settings.UserRepo,
		cache:    settings.Cache,
	}

	return ctrl
}

func (c *DocumentsController) UploadDocument(
	ctx context.Context,
	meta models.Metadata,
	file io.Reader,
) error {
	// Invalidate user cache
	if err := c.cache.InvalidateUserCache(ctx, *meta.OwnerID); err != nil {
		return err
	}

	// If file is binary then upload it to repository
	if meta.File {
		// Upload file to repository
		err := c.fileRepo.UploadFile(ctx, file, meta.FileSize)
		if err != nil {
			return err
		}
	}

	// Save metadata to repository
	err := c.metaRepo.UploadMetadata(ctx, meta)
	if err != nil {
		return err
	}

	return nil
}

func (c *DocumentsController) GetFilesInfo(
	ctx context.Context,
	userID uuid.UUID,
	req models.FilesListRequest,
) ([]models.Metadata, error) {
	// Assign user ID
	id := userID
	if req.Login != "" {
		user, err := c.userRepo.GetUserByLogin(ctx, req.Login)
		if err != nil {
			return nil, err
		}
		id = user.ID
	}

	// Try to get data from cache
	metadata, err := c.cache.GetUserCache(ctx, userID)
	switch {
	case errors.Is(err, apperrors.ErrNotFound):
		// If cache is empty then get data from repository
		metadata, err = c.metaRepo.GetMetadataByUserID(ctx, id)
		if err != nil {
			return nil, err
		}

		// Save data to cache
		err = c.cache.SaveUserCache(ctx, userID, metadata)
		if err != nil {
			return nil, err
		}
	case err != nil:
		return nil, err
	}

	// Create filter
	filter := docfilter.New(req.Limit, req.Offset)
	err = filter.AddFilter(req.Key, req.Value)
	if err != nil {
		return nil, err
	}

	// Filter metadata
	metadata = filter.FilterData(metadata)

	// Return filtered data
	return metadata, nil
}

func (c *DocumentsController) GetFileInfo(
	ctx context.Context,
	docID, userID uuid.UUID,
) (models.Metadata, error) {
	// Try to get metadata from cache
	metadata, err := c.cache.GetUserCache(ctx, userID)
	switch {
	case errors.Is(err, apperrors.ErrNotFound):
		// If cache is empty then get data from repository
		metadata, err = c.metaRepo.GetMetadataByUserID(ctx, userID)
		if err != nil {
			return models.Metadata{}, err
		}

		// Save data to cache
		err = c.cache.SaveUserCache(ctx, userID, metadata)
		if err != nil {
			return models.Metadata{}, err
		}
	case err != nil:
		return models.Metadata{}, err
	}

	// Create filter
	filter := docfilter.New(1, 0)
	err = filter.AddFilter("id", docID.String())
	if err != nil {
		return models.Metadata{}, err
	}

	// Filter metadata
	metadata = filter.FilterData(metadata)

	// Return filtered data
	if len(metadata) > 0 {
		return metadata[0], nil
	}

	return models.Metadata{}, apperrors.ErrNotFound
}

func (c *DocumentsController) GetFile(
	ctx context.Context,
	meta models.Metadata,
) (io.ReadSeekCloser, error) {
	return c.fileRepo.GetFile(ctx, meta)
}

func (c *DocumentsController) DeleteFile(ctx context.Context, id, userID uuid.UUID) error {
	// Invalidate user cache
	if err := c.cache.InvalidateUserCache(ctx, userID); err != nil {
		return err
	}

	// Delete file from repository
	err := c.fileRepo.DeleteFile(ctx, id.String())
	if err != nil {
		return err
	}

	// Delete metadata from repository
	err = c.metaRepo.DeleteMetadata(ctx, id, userID)
	if err != nil {
		return err
	}

	return nil
}

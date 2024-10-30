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

// DocumentController used to upload, download and delete documents.
type FileRepository interface {
	// UploadFile upload file to repository.
	// Returns error if upload failed.
	UploadFile(ctx context.Context, file io.Reader, meta models.Metadata) error

	// GetFile get file from repository.
	// Returns error if get failed.
	// Returns io.ReadSeekCloser if get was successful.
	GetFile(ctx context.Context, meta models.Metadata) (io.ReadSeekCloser, error)

	// DeleteFile delete file from repository.
	// Returns error if delete failed.
	DeleteFile(ctx context.Context, id string) error
}

// MetadataRepository used to upload, download and delete metadata.
type MetadataRepository interface {
	// UploadMetadata upload metadata to repository.
	// Returns error if upload failed.
	// Returns file id if upload was successful.
	UploadMetadata(ctx context.Context, meta models.Metadata) (uuid.UUID, error)

	// GetMetadataByUserID get metadata from repository.
	// Returns error if get failed.
	// Returns []models.Metadata if get was successful.
	GetMetadataByUserID(ctx context.Context, userID uuid.UUID) ([]models.Metadata, error)

	// DeleteMetadata delete metadata from repository.
	// Returns error if delete failed.
	DeleteMetadata(ctx context.Context, id, userID uuid.UUID) error
}

// UserRepository used to get user by login.
type UserRepository interface {
	// GetUserByLogin get user from repository.
	// Returns error if get failed.
	// Returns models.User if get was successful.
	GetUserByLogin(ctx context.Context, login string) (models.User, error)
}

// MetadataCache used to cache metadata.
type MetadataCache interface {
	// InvalidateUserCache invalidate user cache.
	// Returns error if invalidate failed.
	InvalidateUserCache(ctx context.Context, id uuid.UUID) error

	// SaveUserCache save user cache.
	// Returns error if save failed.
	SaveUserCache(ctx context.Context, id uuid.UUID, meta []models.Metadata) error

	// GetUserCache get user cache.
	// Returns error if get failed.
	// Returns []models.Metadata if get was successful.
	GetUserCache(ctx context.Context, id uuid.UUID) ([]models.Metadata, error)
}

// Settings used to create DocumentsController.
// Settings must be provided to New function.
// All fields are required and cant be nil.
type Settings struct {
	// FileRepo used to upload, download and delete files.
	FileRepo FileRepository

	// MetaRepo used to upload, download and delete metadata.
	MetaRepo MetadataRepository

	// UserRepo used to get user by login.
	UserRepo UserRepository

	// Cache used to cache metadata.
	Cache MetadataCache
}

// DocumentsController used to upload, download and delete documents.
// Must be initialized with New function.
type DocumentsController struct {
	fileRepo FileRepository
	metaRepo MetadataRepository
	userRepo UserRepository
	cache    MetadataCache
}

// New creates new DocumentsController.
// Returns pointer to DocumentsController.
// Accepts Settings as argument.
func New(settings Settings) *DocumentsController {
	ctrl := &DocumentsController{
		fileRepo: settings.FileRepo,
		metaRepo: settings.MetaRepo,
		userRepo: settings.UserRepo,
		cache:    settings.Cache,
	}

	return ctrl
}

// UploadDocument upload document to repository.
// Returns error if upload failed.
// Returns nil if upload was successful.
// If meta.File is true, file cant be nil.
// If meta.File is false, meta.JSON must be provided.
func (c *DocumentsController) UploadDocument(
	ctx context.Context,
	meta models.Metadata,
	file io.Reader,
) error {
	// Invalidate user cache
	if err := c.cache.InvalidateUserCache(ctx, *meta.OwnerID); err != nil {
		return err
	}

	// Save metadata to repository
	id, err := c.metaRepo.UploadMetadata(ctx, meta)
	if err != nil {
		return err
	}

	meta.ID = &id

	// If file is binary then upload it to repository
	if meta.File {
		// Upload file to repository
		err = c.fileRepo.UploadFile(ctx, file, meta)
		if err != nil {
			return err
		}
	}

	return nil
}

// GetFilesInfo returns list of documents for given user.
// If req.Login is empty, userID will be used to find files info.
// If req.Login is not empty then it will be used to find user ID.
// If req.Key and req.Value are not empty then they will be used to filter documents.
// If req.Limit or req.Offset are not zero then they will be used to limit and offset documents.
// If cache is empty then it will be filled with data from repository.
// Returns error if get failed.
// Returns []models.Metadata if get was successful.
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

// GetFileInfo get metadata for given document id.
// First try to get data from cache.
// If cache is empty then get data from repository.
// Save data to cache.
// Create filter with id eq docID.
// Filter metadata.
// Return filtered data.
// If filtered data is empty then return ErrNotFound.
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

// GetFile get file from repository.
// Returns error if get failed.
// Returns io.ReadSeekCloser if get was successful.
func (c *DocumentsController) GetFile(
	ctx context.Context,
	meta models.Metadata,
) (io.ReadSeekCloser, error) {
	return c.fileRepo.GetFile(ctx, meta)
}

// DeleteFile delete file and its metadata from repository.
// Returns error if delete failed.
// Returns nil if delete was successful.
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

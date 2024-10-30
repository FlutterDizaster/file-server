package miniorepo

import (
	"context"
	"fmt"
	"io"

	"github.com/FlutterDizaster/file-server/internal/models"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// Settings used to create MinioRepository.
// Endpoint, AccessKey, SecretKey and Bucket are required.
// UseSSL defaults to false.
type Settings struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	Bucket    string
	UseSSL    bool
}

// MinioRepository used to upload and download files.
// Must be initialized with New function.
type MinioRepository struct {
	client *minio.Client
	bucket string
}

// New creates a new MinioRepository instance.
//
// It takes a context.Context, and Settings to initialize the repository.
//
// It creates a MinioRepository using the given settings, checks if the bucket
// exists, and creates the bucket if it doesn't.
//
// It returns the pointer to created MinioRepository and an error.
func New(ctx context.Context, settings Settings) (*MinioRepository, error) {
	repo := &MinioRepository{
		bucket: settings.Bucket,
	}

	// Create minio client
	client, err := minio.New(settings.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(settings.AccessKey, settings.SecretKey, ""),
		Secure: settings.UseSSL,
	})
	if err != nil {
		return nil, err
	}

	repo.client = client

	// Create bucket if not exists
	exists, err := repo.client.BucketExists(ctx, repo.bucket)
	if err != nil {
		return nil, err
	}

	if !exists {
		err = repo.client.MakeBucket(ctx, repo.bucket, minio.MakeBucketOptions{})
		if err != nil {
			return nil, err
		}
	}

	return repo, nil
}

// UploadFile uploads a file to the Minio repository.
//
// It takes an io.Reader representing the file
// to be uploaded, and metadata containing file information like owner ID and file ID.
//
// The file is uploaded to the bucket specified in the repository, using a
// filename composed of the owner ID and file ID.
//
// Returns an error if the upload fails.
func (r MinioRepository) UploadFile(
	ctx context.Context,
	file io.Reader,
	meta models.Metadata,
) error {
	fileName := fmt.Sprintf("%s:%s", meta.OwnerID.String(), meta.ID.String())
	_, err := r.client.PutObject(
		ctx,
		r.bucket,
		fileName,
		file,
		meta.FileSize,
		minio.PutObjectOptions{},
	)
	return err
}

// GetFile get file from repository.
//
// It takes metadata containing file information like owner ID and file ID.
//
// The file is downloaded from the bucket specified in the repository, using a
// filename composed of the owner ID and file ID.
//
// Returns error if get failed.
// Returns io.ReadSeekCloser if get was successful.
func (r MinioRepository) GetFile(
	ctx context.Context,
	meta models.Metadata,
) (io.ReadSeekCloser, error) {
	fileName := fmt.Sprintf("%s:%s", meta.OwnerID.String(), meta.ID.String())
	return r.client.GetObject(ctx, r.bucket, fileName, minio.GetObjectOptions{})
}

// DeleteFile removes a file from the Minio repository.
//
// It takes a context and a string representing the file ID.
//
// The file is removed from the bucket specified in the repository.
//
// Returns an error if the deletion fails.
func (r MinioRepository) DeleteFile(ctx context.Context, id string) error {
	return r.client.RemoveObject(ctx, r.bucket, id, minio.RemoveObjectOptions{})
}

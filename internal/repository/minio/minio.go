package miniorepo

import (
	"context"
	"io"

	"github.com/FlutterDizaster/file-server/internal/models"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Settings struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	Bucket    string
	UseSSL    bool
}

type MinioRepository struct {
	client *minio.Client
	bucket string
}

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

func (r MinioRepository) UploadFile(ctx context.Context, file io.Reader, fileSize int64) error {
	_, err := r.client.PutObject(ctx, r.bucket, "file", file, fileSize, minio.PutObjectOptions{})
	return err
}

func (r MinioRepository) GetFile(
	ctx context.Context,
	meta models.Metadata,
) (io.ReadSeekCloser, error) {
	return r.client.GetObject(ctx, r.bucket, meta.ID.String(), minio.GetObjectOptions{})
}

func (r MinioRepository) DeleteFile(ctx context.Context, id string) error {
	return r.client.RemoveObject(ctx, r.bucket, id, minio.RemoveObjectOptions{})
}

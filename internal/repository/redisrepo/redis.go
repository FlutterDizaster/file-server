package redisrepo

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/FlutterDizaster/file-server/internal/apperrors"
	"github.com/FlutterDizaster/file-server/internal/models"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

const (
	casheKey = "metadata:"
)

// Settings used to create RedisRepository.
// ConnectionString and TTL are required.
// Must be used with New function.
type Settings struct {
	ConnectionString string
	TTL              time.Duration
}

// RedisRepository used to save and get metadata from redis cache.
// Must be initialized with New function.
type RedisRepository struct {
	client  *redis.Client
	ttl     time.Duration
	connStr string
}

// New creates a new RedisRepository instance.
//
// It takes a context.Context, and Settings to initialize the repository.
//
// It creates a RedisRepository using the given settings, parses the connection
// string, creates a redis client, tests the connection, and returns the pointer
// to created RedisRepository and an error.
//
// It returns the pointer to created RedisRepository and an error.
func New(ctx context.Context, settings Settings) (*RedisRepository, error) {
	repo := &RedisRepository{
		connStr: settings.ConnectionString,
		ttl:     settings.TTL,
	}

	// Create redis client
	opt, err := redis.ParseURL(repo.connStr)
	if err != nil {
		slog.Error("Error while parsing connection string", slog.Any("err", err))
		return nil, err
	}

	repo.client = redis.NewClient(opt)

	// Test connection
	_, err = repo.client.Ping(ctx).Result()
	if err != nil {
		slog.Error("Error while connecting to redis", slog.Any("err", err))
		return nil, err
	}

	return repo, nil
}

// InvalidateUserCache removes the cached metadata for a user based on the user ID.
// It constructs the cache key using the user ID and deletes the corresponding entry from the cache.
// Returns an error if the deletion fails.
func (r RedisRepository) InvalidateUserCache(ctx context.Context, id uuid.UUID) error {
	key := casheKey + id.String()

	// Delete data from cache
	return r.client.Del(ctx, key).Err()
}

// SaveUserCache saves the metadata for a user in the Redis cache.
// It constructs the cache key using the user ID, marshals the metadata to JSON,
// and stores it in the cache with a specified TTL.
// Returns an error if marshaling fails or if saving to the cache encounters an error.
func (r RedisRepository) SaveUserCache(
	ctx context.Context,
	id uuid.UUID,
	meta []models.Metadata,
) error {
	key := casheKey + id.String()

	// Marshal data
	data, err := models.Metadatas(meta).MarshalJSON()
	if err != nil {
		return err
	}

	// Save data to cache
	return r.client.Set(ctx, key, string(data), r.ttl).Err()
}

// GetUserCache gets the cached metadata for a user based on the user ID.
// It constructs the cache key using the user ID and retrieves the corresponding entry from the cache.
// If the entry is not found, it returns an apperrors.ErrNotFound error.
// If the retrieval fails, it returns the error.
// If the retrieval is successful, it unmarshals the retrieved data to a slice of models.Metadata and returns it.
func (r RedisRepository) GetUserCache(
	ctx context.Context,
	id uuid.UUID,
) ([]models.Metadata, error) {
	key := casheKey + id.String()

	// Get data from cache
	data, err := r.client.Get(ctx, key).Result()
	switch {
	case errors.Is(err, redis.Nil):
		return nil, apperrors.ErrNotFound
	case err != nil:
		return nil, err
	}

	// Unmarshal data
	metadata := make(models.Metadatas, 0)
	err = metadata.UnmarshalJSON([]byte(data))
	if err != nil {
		return nil, err
	}

	return metadata, nil
}

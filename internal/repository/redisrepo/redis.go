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

type Settings struct {
	ConnectionString string
	TTL              time.Duration
}

type RedisRepository struct {
	client  *redis.Client
	ttl     time.Duration
	connStr string
}

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

func (r RedisRepository) InvalidateUserCache(ctx context.Context, id uuid.UUID) error {
	key := casheKey + id.String()

	// Delete data from cache
	return r.client.Del(ctx, key).Err()
}

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

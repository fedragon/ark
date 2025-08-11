package db

import (
	"context"
	"time"

	"github.com/fedragon/ark/internal/metrics"

	"github.com/redis/go-redis/v9"
)

type redisRepo struct {
	client *redis.Client
}

func NewRepository(client *redis.Client) Repository {
	return &redisRepo{
		client: client,
	}
}

func (r *redisRepo) Close() error {
	return r.client.Close()
}

func (r *redisRepo) Get(ctx context.Context, hash []byte) (*Media, error) {
	now := time.Now()
	defer func() {
		metrics.GetDurationMs.Observe(float64(time.Since(now).Milliseconds()))
	}()

	data, err := r.client.HGetAll(ctx, string(hash)).Result()
	if err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return nil, nil
	}

	createdAt, err := time.Parse(time.RFC3339Nano, data["created_at"])
	if err != nil {
		return nil, err
	}

	var importedAt *time.Time
	if at, ok := data["imported_at"]; ok {
		imported, err := time.Parse(time.RFC3339Nano, at)
		if err != nil {
			return nil, err
		}

		importedAt = &imported
	}

	return &Media{
		Hash:       hash,
		Path:       data["path"],
		CreatedAt:  createdAt,
		ImportedAt: importedAt,
	}, nil
}

func (r *redisRepo) Store(ctx context.Context, media Media) error {
	now := time.Now()
	defer func() {
		metrics.StoreDurationMs.Observe(float64(time.Since(now).Milliseconds()))
	}()

	values := map[string]interface{}{
		"path":       media.Path,
		"created_at": media.CreatedAt.Format(time.RFC3339Nano),
	}

	if media.ImportedAt != nil {
		values["imported_at"] = media.ImportedAt.Format(time.RFC3339Nano)
	}

	return r.client.HSet(ctx, string(media.Hash), values).Err()
}

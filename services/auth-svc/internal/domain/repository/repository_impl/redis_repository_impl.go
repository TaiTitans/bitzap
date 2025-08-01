package repository

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/taititans/bitzap/auth-svc/internal/domain/repository"
	"github.com/taititans/bitzap/auth-svc/internal/util"
)

// redisRepository implements repository.RedisRepository
type redisRepository struct {
	client *redis.Client
	logger util.Logger
}

// NewRedisRepository creates a new Redis repository
func NewRedisRepository(client *redis.Client, logger util.Logger) repository.RedisRepository {
	return &redisRepository{
		client: client,
		logger: logger,
	}
}

// Set sets a key-value pair with expiration
func (r *redisRepository) Set(ctx context.Context, key, value string, expiration time.Duration) error {
	err := r.client.Set(ctx, key, value, expiration).Err()
	if err != nil {
		r.logger.Error("Failed to set Redis key",
			util.String("key", key),
			util.Error(err))
		return err
	}

	r.logger.Info("Successfully set Redis key",
		util.String("key", key),
		util.String("expiration", expiration.String()))

	return nil
}

// Get gets a value by key
func (r *redisRepository) Get(ctx context.Context, key string) (string, error) {
	value, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			// Key doesn't exist
			return "", nil
		}
		r.logger.Error("Failed to get Redis key",
			util.String("key", key),
			util.Error(err))
		return "", err
	}

	r.logger.Info("Successfully got Redis key",
		util.String("key", key))

	return value, nil
}

// Del deletes a key
func (r *redisRepository) Del(ctx context.Context, key string) error {
	err := r.client.Del(ctx, key).Err()
	if err != nil {
		r.logger.Error("Failed to delete Redis key",
			util.String("key", key),
			util.Error(err))
		return err
	}

	r.logger.Info("Successfully deleted Redis key",
		util.String("key", key))

	return nil
}

// Exists checks if a key exists
func (r *redisRepository) Exists(ctx context.Context, key string) (bool, error) {
	exists, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		r.logger.Error("Failed to check Redis key existence",
			util.String("key", key),
			util.Error(err))
		return false, err
	}

	return exists > 0, nil
}

// Expire sets expiration for a key
func (r *redisRepository) Expire(ctx context.Context, key string, expiration time.Duration) error {
	err := r.client.Expire(ctx, key, expiration).Err()
	if err != nil {
		r.logger.Error("Failed to set expiration for Redis key",
			util.String("key", key),
			util.Error(err))
		return err
	}

	r.logger.Info("Successfully set expiration for Redis key",
		util.String("key", key),
		util.String("expiration", expiration.String()))

	return nil
}

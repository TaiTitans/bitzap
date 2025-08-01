package repository

import (
	"context"
	"time"
)

// RedisRepository defines the interface for Redis operations
type RedisRepository interface {
	// Set sets a key-value pair with expiration
	Set(ctx context.Context, key, value string, expiration time.Duration) error

	// Get gets a value by key
	Get(ctx context.Context, key string) (string, error)

	// Del deletes a key
	Del(ctx context.Context, key string) error

	// Exists checks if a key exists
	Exists(ctx context.Context, key string) (bool, error)

	// Expire sets expiration for a key
	Expire(ctx context.Context, key string, expiration time.Duration) error
}

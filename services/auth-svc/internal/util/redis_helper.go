package util

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisHelper helper functions Redis
type RedisHelper struct {
	client *redis.Client
}

// NewRedisHelper
func NewRedisHelper(client *redis.Client) *RedisHelper {
	return &RedisHelper{client: client}
}

// Set key-value with TTL
func (r *RedisHelper) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, key, jsonValue, ttl).Err()
}

// Get value in key
func (r *RedisHelper) Get(ctx context.Context, key string, dest interface{}) error {
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(val), dest)
}

// Delete key
func (r *RedisHelper) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

// Exists key
func (r *RedisHelper) Exists(ctx context.Context, key string) (bool, error) {
	result, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return result > 0, nil
}

// SetNX set key if not exists (atomic)
func (r *RedisHelper) SetNX(ctx context.Context, key string, value interface{}, ttl time.Duration) (bool, error) {
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return false, err
	}
	return r.client.SetNX(ctx, key, jsonValue, ttl).Result()
}

// Incr increase counter
func (r *RedisHelper) Incr(ctx context.Context, key string) (int64, error) {
	return r.client.Incr(ctx, key).Result()
}

// IncrBy increase counter by value
func (r *RedisHelper) IncrBy(ctx context.Context, key string, value int64) (int64, error) {
	return r.client.IncrBy(ctx, key, value).Result()
}

// Expire set TTL key
func (r *RedisHelper) Expire(ctx context.Context, key string, ttl time.Duration) error {
	return r.client.Expire(ctx, key, ttl).Err()
}

// TTL get TTL key
func (r *RedisHelper) TTL(ctx context.Context, key string) (time.Duration, error) {
	return r.client.TTL(ctx, key).Result()
}

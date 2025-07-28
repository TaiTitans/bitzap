package initialize

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisConfig
type RedisConfig struct {
	Address         string
	Password        string
	DB              int
	IdleTimeout     string
	MaxConnLifetime string
	WaitTimeout     string
	DialTimeout     string
	ReadTimeout     string
	WriteTimeout    string
	MaxActive       int
}

// InitRedis connection pool
func InitRedis(redisConfig RedisConfig) *redis.Client {
	// Parse duration strings
	dialTimeout, _ := time.ParseDuration(redisConfig.DialTimeout)
	readTimeout, _ := time.ParseDuration(redisConfig.ReadTimeout)
	writeTimeout, _ := time.ParseDuration(redisConfig.WriteTimeout)

	rdb := redis.NewClient(&redis.Options{
		Addr:         redisConfig.Address,
		Password:     redisConfig.Password,
		DB:           redisConfig.DB,
		DialTimeout:  dialTimeout,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		PoolSize:     redisConfig.MaxActive,
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	log.Printf("Successfully connected to Redis at %s", redisConfig.Address)
	return rdb
}

// CloseRedis Redis
func CloseRedis(rdb *redis.Client) {
	if rdb != nil {
		if err := rdb.Close(); err != nil {
			log.Printf("Error closing Redis connection: %v", err)
		}
	}
}

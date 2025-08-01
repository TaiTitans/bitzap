package config

import (
	"os"
	"strconv"
)

// Config holds all application configuration
type Config struct {
	Email    EmailConfig
	Database DatabaseConfig
	Redis    RedisConfig
	Server   ServerConfig
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host            string
	Port            string
	User            string
	Password        string
	DBName          string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime int // in hours
}

// RedisConfig holds Redis configuration
type RedisConfig struct {
	Address      string
	Password     string
	DB           int
	DialTimeout  string
	ReadTimeout  string
	WriteTimeout string
	MaxActive    int
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port     string
	LogLevel string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	return &Config{
		Email:    loadEmailConfig(),
		Database: loadDatabaseConfig(),
		Redis:    loadRedisConfig(),
		Server:   loadServerConfig(),
	}
}

// loadEmailConfig loads email configuration from environment variables
func loadEmailConfig() EmailConfig {
	return EmailConfig{
		MailjetAPIKey:    getEnv("MAILJET_API_KEY", ""),
		MailjetSecretKey: getEnv("MAILJET_SECRET_KEY", ""),
		FromEmail:        getEnv("FROM_EMAIL", "noreply@bitzap.com"),
		FromName:         getEnv("FROM_NAME", "Bitzap Auth Service"),
		AppURL:           getEnv("APP_URL", "http://localhost:8080"),
	}
}

// loadDatabaseConfig loads database configuration from environment variables
func loadDatabaseConfig() DatabaseConfig {
	maxOpenConns, _ := strconv.Atoi(getEnv("DB_MAX_OPEN_CONNS", "150"))
	maxIdleConns, _ := strconv.Atoi(getEnv("DB_MAX_IDLE_CONNS", "10"))
	connMaxLifetime, _ := strconv.Atoi(getEnv("DB_CONN_MAX_LIFETIME", "1"))

	return DatabaseConfig{
		Host:            getEnv("DB_HOST", "localhost"),
		Port:            getEnv("DB_PORT", "5432"),
		User:            getEnv("DB_USER", "admin"),
		Password:        getEnv("DB_PASSWORD", "admin123"),
		DBName:          getEnv("DB_NAME", "auth_service"),
		SSLMode:         getEnv("DB_SSL_MODE", "disable"),
		MaxOpenConns:    maxOpenConns,
		MaxIdleConns:    maxIdleConns,
		ConnMaxLifetime: connMaxLifetime,
	}
}

// loadRedisConfig loads Redis configuration from environment variables
func loadRedisConfig() RedisConfig {
	db, _ := strconv.Atoi(getEnv("REDIS_DB", "0"))
	maxActive, _ := strconv.Atoi(getEnv("REDIS_MAX_ACTIVE", "10000"))

	return RedisConfig{
		Address:      getEnv("REDIS_ADDRESS", "127.0.0.1:6379"),
		Password:     getEnv("REDIS_PASSWORD", "redispass"),
		DB:           db,
		DialTimeout:  getEnv("REDIS_DIAL_TIMEOUT", "30s"),
		ReadTimeout:  getEnv("REDIS_READ_TIMEOUT", "30s"),
		WriteTimeout: getEnv("REDIS_WRITE_TIMEOUT", "30s"),
		MaxActive:    maxActive,
	}
}

// loadServerConfig loads server configuration from environment variables
func loadServerConfig() ServerConfig {
	return ServerConfig{
		Port:     getEnv("SERVER_PORT", "8080"),
		LogLevel: getEnv("LOG_LEVEL", "info"),
	}
}

// getEnv gets environment variable with fallback default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

package config

import (
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

// Config holds all application configuration
type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Logger   LoggerConfig   `yaml:"logger"`
	Database DatabaseConfig `yaml:"database"`
	Redis    RedisConfig    `yaml:"redis"`
	Auth     AuthConfig     `yaml:"auth"`
	Email    EmailConfig    `yaml:"email"`
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Address         string `yaml:"address"`
	Port            string `yaml:"port"`
	OpenAPIPath     string `yaml:"openapiPath"`
	SwaggerPath     string `yaml:"swaggerPath"`
	ErrorStack      bool   `yaml:"errorStack"`
	ErrorLogEnabled bool   `yaml:"errorLogEnabled"`
	ErrorLogPattern string `yaml:"errorLogPattern"`
	LogLevel        string // Derived from logger config
}

// LoggerConfig holds logger configuration
type LoggerConfig struct {
	Path                 string   `yaml:"path"`
	File                 string   `yaml:"file"`
	Prefix               string   `yaml:"prefix"`
	Level                string   `yaml:"level"`
	TimeFormat           string   `yaml:"timeFormat"`
	CtxKeys              []string `yaml:"ctxKeys"`
	Header               bool     `yaml:"header"`
	StSkip               int      `yaml:"stSkip"`
	Stdout               bool     `yaml:"stdout"`
	RotateSize           int      `yaml:"rotateSize"`
	RotateExpire         int      `yaml:"rotateExpire"`
	RotateBackupLimit    int      `yaml:"rotateBackupLimit"`
	RotateBackupExpire   int      `yaml:"rotateBackupExpire"`
	RotateBackupCompress int      `yaml:"rotateBackupCompress"`
	RotateCheckInterval  string   `yaml:"rotateCheckInterval"`
	StdoutColorDisabled  bool     `yaml:"stdoutColorDisabled"`
	WriterColorEnable    bool     `yaml:"writerColorEnable"`
	Flags                int      `yaml:"flags"`
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host            string               `yaml:"host"`
	Port            string               `yaml:"port"`
	User            string               `yaml:"user"`
	Password        string               `yaml:"password"`
	DBName          string               `yaml:"dbname"`
	SSLMode         string               `yaml:"sslmode"`
	MaxOpenConns    int                  `yaml:"max_open_conns"`
	MaxIdleConns    int                  `yaml:"max_idle_conns"`
	ConnMaxLifetime int                  `yaml:"conn_max_lifetime"`
	Logger          DatabaseLoggerConfig `yaml:"logger"`
	Debug           bool                 `yaml:"debug"`
}

type DatabaseLoggerConfig struct {
	Level  string `yaml:"level"`
	Stdout bool   `yaml:"stdout"`
}

// RedisConfig holds Redis configuration
type RedisConfig struct {
	Default RedisDefaultConfig `yaml:"default"`
}

type RedisDefaultConfig struct {
	Address         string `yaml:"address"`
	Password        string `yaml:"password"`
	DB              int    `yaml:"db"`
	IdleTimeout     string `yaml:"idleTimeout"`
	MaxConnLifetime string `yaml:"maxConnLifetime"`
	WaitTimeout     string `yaml:"waitTimeout"`
	DialTimeout     string `yaml:"dialTimeout"`
	ReadTimeout     string `yaml:"readTimeout"`
	WriteTimeout    string `yaml:"writeTimeout"`
	MaxActive       int    `yaml:"maxActive"`
}

// AuthConfig holds authentication configuration
type AuthConfig struct {
	SecretKey                string `yaml:"secretKey"`
	AccessTokenExpireMinute  int    `yaml:"accessTokenExpireMinute"`
	RefreshTokenExpireMinute int    `yaml:"refreshTokenExpireMinute"`
}

// LoadConfig loads configuration from YAML file
func LoadConfig() *Config {
	data, err := ioutil.ReadFile("configs/config.yaml")
	if err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		log.Fatalf("Failed to parse config file: %v", err)
	}

	// Set derived values
	config.Server.LogLevel = config.Logger.Level
	if config.Server.Port == "" {
		// Extract port from address if not set separately
		if config.Server.Address != "" {
			config.Server.Port = config.Server.Address[1:]
		}
	}

	// Load sensitive data from environment variables
	config.Email.MailjetAPIKey = getEnv("MAILJET_API_KEY", config.Email.MailjetAPIKey)
	config.Email.MailjetSecretKey = getEnv("MAILJET_SECRET_KEY", config.Email.MailjetSecretKey)

	return &config
}

// getEnv gets environment variable with fallback default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

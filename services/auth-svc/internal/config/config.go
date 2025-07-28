package config

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config
type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Logger   LoggerConfig   `yaml:"logger"`
	Database DatabaseConfig `yaml:"database"`
	Redis    RedisConfig    `yaml:"redis"`
	Auth     AuthConfig     `yaml:"auth"`
}

// ServerConfig
type ServerConfig struct {
	Address         string `yaml:"address"`
	OpenAPIPath     string `yaml:"openapiPath"`
	SwaggerPath     string `yaml:"swaggerPath"`
	ErrorStack      bool   `yaml:"errorStack"`
	ErrorLogEnabled bool   `yaml:"errorLogEnabled"`
	ErrorLogPattern string `yaml:"errorLogPattern"`
}

// LoggerConfig
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

// DatabaseConfig
type DatabaseConfig struct {
	Logger  DatabaseLoggerConfig  `yaml:"logger"`
	Default DatabaseDefaultConfig `yaml:"default"`
}

type DatabaseLoggerConfig struct {
	Level  string `yaml:"level"`
	Stdout bool   `yaml:"stdout"`
}

type DatabaseDefaultConfig struct {
	Link        string `yaml:"link"`
	Debug       bool   `yaml:"debug"`
	MaxIdle     string `yaml:"maxIdle"`
	MaxOpen     string `yaml:"maxOpen"`
	MaxLifetime string `yaml:"maxLifetime"`
}

// RedisConfig
type RedisConfig struct {
	Default RedisDefaultConfig `yaml:"default"`
}

type RedisDefaultConfig struct {
	Address         string `yaml:"address"`
	DB              int    `yaml:"db"`
	IdleTimeout     string `yaml:"idleTimeout"`
	MaxConnLifetime string `yaml:"maxConnLifetime"`
	WaitTimeout     string `yaml:"waitTimeout"`
	DialTimeout     string `yaml:"dialTimeout"`
	ReadTimeout     string `yaml:"readTimeout"`
	WriteTimeout    string `yaml:"writeTimeout"`
	MaxActive       int    `yaml:"maxActive"`
}

// AuthConfig
type AuthConfig struct {
	SecretKey                string `yaml:"secretKey"`
	AccessTokenExpireMinute  int    `yaml:"accessTokenExpireMinute"`
	RefreshTokenExpireMinute int    `yaml:"refreshTokenExpireMinute"`
}

// LoadConfig
func LoadConfig(configPath string) (*Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

// GetRedisConfig
func (c *Config) GetRedisConfig() map[string]interface{} {
	redisConfig := make(map[string]interface{})

	// Parse duration strings
	idleTimeout, _ := time.ParseDuration(c.Redis.Default.IdleTimeout)
	maxConnLifetime, _ := time.ParseDuration(c.Redis.Default.MaxConnLifetime)
	waitTimeout, _ := time.ParseDuration(c.Redis.Default.WaitTimeout)
	dialTimeout, _ := time.ParseDuration(c.Redis.Default.DialTimeout)
	readTimeout, _ := time.ParseDuration(c.Redis.Default.ReadTimeout)
	writeTimeout, _ := time.ParseDuration(c.Redis.Default.WriteTimeout)

	redisConfig["address"] = c.Redis.Default.Address
	redisConfig["db"] = c.Redis.Default.DB
	redisConfig["idleTimeout"] = idleTimeout
	redisConfig["maxConnLifetime"] = maxConnLifetime
	redisConfig["waitTimeout"] = waitTimeout
	redisConfig["dialTimeout"] = dialTimeout
	redisConfig["readTimeout"] = readTimeout
	redisConfig["writeTimeout"] = writeTimeout
	redisConfig["maxActive"] = c.Redis.Default.MaxActive

	return redisConfig
}

package initialize

import (
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

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

// InitLogger Zap logger
func InitLogger(config LoggerConfig) *zap.Logger {
	if err := os.MkdirAll(config.Path, 0755); err != nil {
		panic("Failed to create log directory: " + err.Error())
	}

	// Parse log level
	level, err := zapcore.ParseLevel(config.Level)
	if err != nil {
		level = zapcore.InfoLevel
	}

	// encoder config
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder

	// cores
	var cores []zapcore.Core

	// File core rotation
	if config.File != "" {
		filePath := filepath.Join(config.Path, config.File)
		fileCore := zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			zapcore.AddSync(&lumberjack.Logger{
				Filename:   filePath,
				MaxSize:    config.RotateSize, // MB
				MaxBackups: config.RotateBackupLimit,
				MaxAge:     config.RotateExpire, // days
				Compress:   config.RotateBackupCompress > 0,
			}),
			level,
		)
		cores = append(cores, fileCore)
	}

	// Stdout core
	if config.Stdout {
		stdoutCore := zapcore.NewCore(
			zapcore.NewConsoleEncoder(encoderConfig),
			zapcore.AddSync(os.Stdout),
			level,
		)
		cores = append(cores, stdoutCore)
	}

	// logger
	core := zapcore.NewTee(cores...)
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(config.StSkip))

	return logger
}

// InitSugarLogger
func InitSugarLogger(config LoggerConfig) *zap.SugaredLogger {
	return InitLogger(config).Sugar()
}

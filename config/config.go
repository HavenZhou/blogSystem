package config

import (
	"errors"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	DB struct {
		DSN         string
		MaxIdleConn int
		MaxOpenConn int
	}
	JWT struct {
		Secret   string
		Lifetime time.Duration
	}
	Server struct {
		Port string
		Env  string
	}
	Log struct {
		Level string
	}
}

func Load() (*Config, error) {
	_ = godotenv.Load(".env", ".env.local")

	cfg := &Config{
		DB: struct {
			DSN         string
			MaxIdleConn int
			MaxOpenConn int
		}{
			DSN:         getEnv("DB_DSN", ""),
			MaxIdleConn: 10,
			MaxOpenConn: 100,
		},
		JWT: struct {
			Secret   string
			Lifetime time.Duration
		}{
			Secret:   getEnv("JWT_SECRET", ""),
			Lifetime: 24 * time.Hour,
		},
		Server: struct {
			Port string
			Env  string
		}{
			Port: getEnv("SERVER_PORT", "8080"),
			Env:  getEnv("SERVER_ENV", "development"),
		},
		Log: struct {
			Level string
		}{
			Level: getEnv("LOG_LEVEL", "info"),
		},
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}
	return cfg, nil
}

func (c *Config) validate() error {
	// 验证数据库配置
	if c.DB.DSN == "" {
		return errors.New("database DSN configuration is required")
	}

	// 验证JWT配置
	if c.JWT.Secret == "" {
		return errors.New("JWT secret key is required")
	}
	if len(c.JWT.Secret) < 32 {
		return errors.New("JWT secret must be at least 32 characters long")
	}

	// 验证服务器配置
	if c.Server.Port == "" {
		return errors.New("server port is required")
	}
	if _, err := strconv.Atoi(c.Server.Port); err != nil {
		return errors.New("server port must be a valid number")
	}
	if c.Server.Env != "development" && c.Server.Env != "production" {
		return errors.New("server environment must be either 'development' or 'production'")
	}

	// 验证日志级别
	validLogLevels := map[string]bool{
		"debug":  true,
		"info":   true,
		"warn":   true,
		"error":  true,
		"dpanic": true,
		"panic":  true,
		"fatal":  true,
	}
	if !validLogLevels[strings.ToLower(c.Log.Level)] {
		return errors.New("invalid log level, must be one of: debug, info, warn, error, dpanic, panic, fatal")
	}

	return nil
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

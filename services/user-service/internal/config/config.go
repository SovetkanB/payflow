package config

import (
	"fmt"
	"os"
	"time"
)

type Config struct {
	HTTPPort      string
	DBHost        string
	DBPort        string
	DBUser        string
	DBPass        string
	DBName        string
	DBSSLMode     string
	JWTSecret     string
	JWTExpiretion time.Duration
}

func (c *Config) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.DBHost, c.DBPort, c.DBUser, c.DBPass, c.DBName, c.DBSSLMode,
	)
}

func Load() (*Config, error) {
	jwtExpiretion, err := time.ParseDuration(getEnv("JWT_EXPIRETION", "15m"))
	if err != nil {
		return nil, fmt.Errorf("invalid JWT_ACCESS_TTL: %w", err)
	}
	return &Config{
		HTTPPort:      getEnv("APP_PORT", "8080"),
		DBHost:        getEnv("DB_HOST", "localhost"),
		DBPort:        getEnv("DB_PORT", "5432"),
		DBUser:        getEnv("DB_USER", "postgres"),
		DBPass:        getEnv("DB_PASS", "postgres"),
		DBName:        getEnv("DB_NAME", "userdb"),
		DBSSLMode:     getEnv("DB_SSL_MODE", "disable"),
		JWTSecret:     getEnv("JWT_SECRET", "changeinprod"),
		JWTExpiretion: jwtExpiretion,
	}, nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}

	return fallback
}

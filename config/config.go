package config

import (
	"fmt"
	"sync"

	"github.com/jessevdk/go-flags"
	"github.com/joho/godotenv"
)

var (
	instance *Config
	once     sync.Once
)

// LoadConfig loads the configuration from .env and command-line flags.
func LoadConfig() (*Config, error) {
	var cfg Config
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}
	fp := flags.NewParser(&cfg, flags.Default)
	// Parse flags
	if _, err := fp.Parse(); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// GetConfig returns the singleton instance of Config.
func GetConfig() (*Config, error) {
	var err error
	once.Do(func() {
		instance, err = LoadConfig()
	})
	return instance, err
}

type Config struct {
	Command struct {
		Migrate bool `long:"migrate"`
		Seed    bool `long:"seed"`
	}

	Server struct {
		Port int `long:"server-port" env:"SERVER_PORT" default:"8080"`
	}

	JWT struct {
		PrivateKey string `long:"jwt-private-key" env:"JWT_PRIVATE_KEY"`
		Duration   struct {
			AccessTokenInMin   int `long:"jwt-duration-access-token-in-min" env:"JWT_DURATION_ACCESS_TOKEN_IN_MIN" default:"15"`
			RefreshTokenInHour int `long:"jwt-duration-refresh-token-in-hour" env:"JWT_DURATION_REFRESH_TOKEN_IN_HOUR" default:"24"`
			KeepLogInDay       int `long:"jwt-duration-keep-log-in-day" env:"JWT_DURATION_KEEP_LOG_IN_DAY" default:"30"`
		}
	}

	Postgres struct {
		DBName   string `long:"postgres-db-name" env:"POSTGRES_DB_NAME" default:"federico"`
		User     string `long:"postgres-user" env:"POSTGRES_USER" default:"postgres"`
		Password string `long:"postgres-password" env:"POSTGRES_PASSWORD"`
		Host     string `long:"postgres-host" env:"POSTGRES_HOST" default:"localhost"`
		Port     int    `long:"postgres-port" env:"POSTGRES_PORT" default:"5432"`
		SSLMode  string `long:"postgres-ssl-mode" env:"POSTGRES_SSL_MODE" default:"disable"`
	}
}

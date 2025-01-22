package config

import (
	"fmt"
	"os"
	"strconv"
	"sync"

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

	// fp := flags.NewParser(&cfg, flags.Default)
	// // Parse flags
	// if _, err := fp.Parse(); err != nil {
	// 	return nil, err
	// }

	serverPort, err := strconv.Atoi(os.Getenv("SERVER_PORT"))
	if err != nil {
		cfg.Server.Port = 8080
	} else {
		cfg.Server.Port = serverPort
	}

	cfg.Postgres.Host = os.Getenv("POSTGRES_HOST")
	postgresPort, err := strconv.Atoi(os.Getenv("POSTGRES_PORT"))
	if err != nil {
		cfg.Postgres.Port = 5432
	} else {
		cfg.Postgres.Port = postgresPort
	}
	cfg.Postgres.DBName = os.Getenv("POSTGRES_DB_NAME")
	cfg.Postgres.User = os.Getenv("POSTGRES_USER")
	cfg.Postgres.Password = os.Getenv("POSTGRES_PASSWORD")
	cfg.Postgres.SSLMode = os.Getenv("POSTGRES_SSL_MODE")

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

	Postgres struct {
		DBName   string `long:"postgres-db-name" env:"POSTGRES_DB_NAME" default:"federico"`
		User     string `long:"postgres-user" env:"POSTGRES_USER" default:"postgres"`
		Password string `long:"postgres-password" env:"POSTGRES_PASSWORD"`
		Host     string `long:"postgres-host" env:"POSTGRES_HOST" default:"localhost"`
		Port     int    `long:"postgres-port" env:"POSTGRES_PORT" default:"5432"`
		SSLMode  string `long:"postgres-ssl-mode" env:"POSTGRES_SSL_MODE" default:"disable"`
	}
}

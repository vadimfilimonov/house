package config

import (
	"fmt"
	"log"

	env "github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

const (
	defaultServerAddress = "localhost:8080"
)

type Config struct {
	// HTTP server startup address
	ServerAddress string `env:"SERVER_ADDRESS"`
	// Address of the database connection
	DatabaseURL string `env:"DATABASE_URL"`
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file")
	}
}

func New() Config {
	return Config{}
}

func (c *Config) Parse() error {
	err := env.Parse(c)
	if err != nil {
		return fmt.Errorf("cannot parse env: %w", err)
	}

	if c.ServerAddress == "" {
		c.ServerAddress = defaultServerAddress
	}

	return nil
}

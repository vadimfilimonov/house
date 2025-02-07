package config

import (
	"fmt"
	"log"

	env "github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

const (
	defaultServerAddress = "localhost:8080"
	defaultRedisAddress  = "localhost:6379"
)

type Config struct {
	// HTTP server startup address
	ServerAddress string `env:"SERVER_ADDRESS"`
	// Address of the database connection
	DatabaseURL string `env:"DATABASE_URL"`
	// Address of the redis connection
	RedisAddress string `env:"REDIS_ADDRESS"`
	// Password of the redis
	RedisPassword string `env:"REDIS_PASSWORD"`
	// The secret key for signing the JWT token
	JwtSecretKey string `env:"JWT_SECRET_KEY"`
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
	if err := env.Parse(c); err != nil {
		return fmt.Errorf("cannot parse env: %w", err)
	}

	if c.ServerAddress == "" {
		c.ServerAddress = defaultServerAddress
	}

	if c.RedisAddress == "" {
		c.RedisAddress = defaultRedisAddress
	}

	return nil
}

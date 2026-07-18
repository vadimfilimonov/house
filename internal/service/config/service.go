package config

import (
	"fmt"
	"os"

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
	// Host of the database
	DatabaseHost string `env:"DB_HOST"`
	// Port of the database
	DatabasePort string `env:"DB_PORT"`
	// Postgres user
	PostgresUser string `env:"POSTGRES_USER"`
	// Postgres password
	PostgresPassword string `env:"POSTGRES_PASSWORD"`
	// Postgres database name
	PostgresDatabaseName string `env:"POSTGRES_DB"`
	// Host for the redis connection
	RedisHost string `env:"REDIS_HOST"`
	// Port for the redis connection
	RedisPort string `env:"REDIS_PORT"`
	// Password of the redis
	RedisPassword string `env:"REDIS_PASSWORD"`
	// The secret key for signing the JWT token
	JwtSecretKey string `env:"JWT_SECRET_KEY"`
}

func New() Config {
	return Config{}
}

func (c *Config) Parse() error {
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("cannot load .env file: %w", err)
	}

	if err := env.Parse(c); err != nil {
		return fmt.Errorf("cannot parse env: %w", err)
	}

	if c.ServerAddress == "" {
		c.ServerAddress = defaultServerAddress
	}

	if c.RedisHost == "" {
		c.RedisHost = defaultRedisAddress
	}

	return nil
}

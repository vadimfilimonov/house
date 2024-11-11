package token

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	tokenKeyPrefix = "token"
)

type Database struct {
	db *redis.Client
}

func New(address, password string) (*Database, error) {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       0,
	})

	_, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	return &Database{
		db: redisClient,
	}, nil
}

func (d *Database) Add(
	ctx context.Context,
	key string,
	value string,
	expiration time.Duration,
) error {
	statusCmd := d.db.Set(ctx, buildKey(key), value, expiration)
	if statusCmd == nil {
		return fmt.Errorf("redis status cmd is nil")
	}

	if err := statusCmd.Err(); err != nil {
		return err
	}

	return nil
}

func (d *Database) Get(ctx context.Context, key string) (string, error) {
	stringCmd := d.db.Get(ctx, buildKey(key))
	if stringCmd == nil {
		return "", fmt.Errorf("redis string cmd is nil")
	}

	return stringCmd.Result()
}

func (d *Database) Has(ctx context.Context, key string) bool {
	_, err := d.Get(ctx, buildKey(key))
	return err == nil
}

func (d *Database) Close() error {
	return d.db.Close()
}

func buildKey(key string) string {
	return fmt.Sprintf("%s:%s", tokenKeyPrefix, key)
}

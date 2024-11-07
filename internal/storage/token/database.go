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

func NewDatabase(address, password string) (*Database, error) {
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
	key, value string,
	expiration time.Duration,
) error {
	err := d.db.Set(ctx, buildKey(key), value, expiration).Err()
	if err != nil {
		return err
	}

	return nil
}

func (d *Database) Get(ctx context.Context, key string) (string, error) {
	token, err := d.db.Get(ctx, buildKey(key)).Result()
	if err != nil {
		return "", err
	}

	return token, nil
}

func buildKey(key string) string {
	return fmt.Sprintf("%s:%s", tokenKeyPrefix, key)
}

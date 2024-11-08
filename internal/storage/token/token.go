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
	expiration time.Duration,
) error {
	err := d.db.Set(ctx, buildKey(key), "", expiration).Err()
	if err != nil {
		return err
	}

	return nil
}

func (d *Database) Has(ctx context.Context, key string) bool {
	_, err := d.db.Get(ctx, buildKey(key)).Result()
	return err == nil
}

func (d *Database) Close() error {
	return d.db.Close()
}

func buildKey(key string) string {
	return fmt.Sprintf("%s:%s", tokenKeyPrefix, key)
}

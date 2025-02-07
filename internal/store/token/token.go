package token

import (
	"context"
	"fmt"
	"time"

	"github.com/vadimfilimonov/house/internal/storage/redis"
)

const (
	tokenKeyPrefix = "token"
)

type Token struct {
	storage *redis.Storage
}

func New(storage *redis.Storage) *Token {
	return &Token{storage: storage}
}

func (t *Token) Add(
	ctx context.Context,
	key string,
	value string,
	expiration time.Duration,
) error {
	statusCmd := t.storage.Set(ctx, buildKey(key), value, expiration)
	if statusCmd == nil {
		return fmt.Errorf("redis status cmd is nil")
	}

	if err := statusCmd.Err(); err != nil {
		return err
	}

	return nil
}

func (t *Token) Get(ctx context.Context, key string) (string, error) {
	return t.storage.Get(ctx, buildKey(key))
}

func (t *Token) Has(ctx context.Context, key string) bool {
	return t.storage.Has(ctx, buildKey(key))
}

func buildKey(key string) string {
	return fmt.Sprintf("%s:%s", tokenKeyPrefix, key)
}

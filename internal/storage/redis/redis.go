package redis

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

type Storage struct {
	db *redis.Client
}

var (
	storage  = &Storage{}
	syncErr  error
	initOnce = sync.Once{}
)

func New(ctx context.Context, address, password string) (*Storage, error) {
	initOnce.Do(func() {
		client := redis.NewClient(&redis.Options{
			Addr:     address,
			Password: password,
			DB:       0,
		})

		statusCmd := client.Ping(ctx)
		if statusCmd == nil {
			syncErr = fmt.Errorf("status cmd is nil")
			return
		}

		if _, err := statusCmd.Result(); err != nil {
			syncErr = err
			return
		}

		storage.db = client
	})

	return storage, syncErr
}

func (s *Storage) Get(ctx context.Context, key string) (string, error) {
	stringCmd := s.db.Get(ctx, key)
	if stringCmd == nil {
		return "", fmt.Errorf("redis string cmd is nil")
	}

	return stringCmd.Result()
}

func (s *Storage) Set(ctx context.Context, key string, value any, expiration time.Duration) *redis.StatusCmd {
	return s.db.Set(ctx, key, value, expiration)
}

func (s *Storage) Has(ctx context.Context, key string) bool {
	_, err := s.Get(ctx, key)
	return err == nil
}

func (s *Storage) Close() error {
	return s.db.Close()
}

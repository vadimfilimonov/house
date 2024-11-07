package token

import (
	"context"
	"fmt"
	"time"

	"github.com/vadimfilimonov/house/internal/storage"
)

type Storage interface {
	Add(ctx context.Context, key, value string, expiration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
}

func GetStorage(storageType, address, password string) (Storage, error) {
	if storageType == storage.StorageTypeDatabase {
		return NewDatabase(address, password)
	}

	return nil, fmt.Errorf("unknown storage type: %s", storageType)
}

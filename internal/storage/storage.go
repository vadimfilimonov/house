package storage

import (
	"crypto/rand"
	"fmt"

	"github.com/vadimfilimonov/house/internal/models"
)

type Storage interface {
	Add(email, hashedPassword, userType string) (id *string, err error)
	Get(email string) (*models.User, error)
}

func GetStorage(storageType, databaseURL string) (Storage, error) {
	if storageType == StorageTypeMemory {
		return NewMemory(), nil
	}

	if storageType == StorageTypeDatabase {
		db, err := NewDatabase(databaseURL)
		if err != nil {
			return nil, fmt.Errorf("cannot return database storage: %w", err)
		}

		return db, nil
	}

	return nil, fmt.Errorf("unknown storage type: %s", storageType)
}

func generateUserID() (*string, error) {
	b := make([]byte, 16)

	_, err := rand.Read(b)
	if err != nil {
		return nil, fmt.Errorf("cannot generate userID: %s", err.Error())
	}

	uuid := fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	return &uuid, nil
}

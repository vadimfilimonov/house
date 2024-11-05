package user

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"github.com/vadimfilimonov/house/internal/models"
)

type storage interface {
	Add(email, hashedPassword, userType string) (id *string, err error)
	Get(email string) (*models.User, error)
}

type UserManager struct {
	storage storage
}

func New(storage storage) *UserManager {
	return &UserManager{
		storage: storage,
	}
}

func (u *UserManager) Register(email, password, userType string) (*string, error) {
	err := u.validate(email, password, userType)
	if err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	hash, err := hashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("cannot hash password: %w", err)
	}

	id, err := u.storage.Add(email, hash, userType)
	if err != nil {
		return nil, fmt.Errorf("cannot save user in storage: %w", err)
	}

	return id, nil
}

// hashPassword generates a bcrypt hash for the given password.
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

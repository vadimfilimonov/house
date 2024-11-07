package user

import (
	"context"
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/vadimfilimonov/house/internal/models"
)

const (
	WrongPasswordErr = "password is wrong"
)

type userStorage interface {
	Add(email, hashedPassword, userType string) (id *string, err error)
	Get(id string) (*models.User, error)
}

type tokenStorage interface {
	Add(ctx context.Context, key, value string, expiration time.Duration) error
}

type tokenManager interface {
	Encode(sub string) (*string, error)
}

type UserManager struct {
	userStorage  userStorage
	tokenStorage tokenStorage
	tokenManager tokenManager
}

func New(userStorage userStorage, tokenStorage tokenStorage, tokenManager tokenManager) *UserManager {
	return &UserManager{
		userStorage:  userStorage,
		tokenStorage: tokenStorage,
		tokenManager: tokenManager,
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

	id, err := u.userStorage.Add(email, hash, userType)
	if err != nil {
		return nil, fmt.Errorf("cannot save user in storage: %w", err)
	}

	return id, nil
}

func (u *UserManager) Login(ctx context.Context, id, password string) (*string, error) {
	user, err := u.userStorage.Get(id)
	if err != nil {
		return nil, err
	}

	isPasswordCorrect := verifyPassword(password, user.Password)
	if !isPasswordCorrect {
		return nil, errors.New(WrongPasswordErr)
	}

	token, err := u.tokenManager.Encode(user.ID)
	if err != nil {
		return nil, err
	}

	err = u.tokenStorage.Add(ctx, user.ID, *token, 24*time.Hour)
	if err != nil {
		return nil, err
	}

	return token, nil
}

// hashPassword generates a bcrypt hash for the given password.
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// VerifyPassword verifies if the given password matches the stored hash.
func verifyPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

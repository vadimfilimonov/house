package user

import (
	"context"
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/vadimfilimonov/house/internal/models"
)

var (
	ErrWrongPassword = errors.New("password is wrong")
)

type userStore interface {
	Add(ctx context.Context, email, hashedPassword, userType string) (id *string, err error)
	Get(ctx context.Context, id string) (*models.User, error)
}

type tokenStore interface {
	Add(ctx context.Context, key, value string, expiration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
}

type tokenManager interface {
	Encode(sub, userType string) (*string, error)
}

type UserManager struct {
	userStore    userStore
	tokenStore   tokenStore
	tokenManager tokenManager
}

func New(userStorage userStore, tokenStorage tokenStore, tokenManager tokenManager) *UserManager {
	return &UserManager{
		userStore:    userStorage,
		tokenStore:   tokenStorage,
		tokenManager: tokenManager,
	}
}

func (u *UserManager) Register(ctx context.Context, email, password, userType string) (*string, error) {
	if err := u.validate(email, password, userType); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	hash, err := hashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("cannot hash password: %w", err)
	}

	id, err := u.userStore.Add(ctx, email, hash, userType)
	if err != nil {
		return nil, fmt.Errorf("cannot save user in storage: %w", err)
	}

	return id, nil
}

func (u *UserManager) Login(ctx context.Context, id, password string) (*string, error) {
	user, err := u.userStore.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	isPasswordCorrect := verifyPassword(password, user.Password)
	if !isPasswordCorrect {
		return nil, ErrWrongPassword
	}

	savedToken, err := u.tokenStore.Get(ctx, user.ID)
	if err == nil {
		return &savedToken, nil
	}

	token, err := u.tokenManager.Encode(user.ID, user.UserType)
	if err != nil {
		return nil, err
	} else if token == nil {
		return nil, fmt.Errorf("token is nil")
	}

	if err := u.tokenStore.Add(ctx, user.ID, *token, 24*time.Hour); err != nil {
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

package user

import (
	"context"
	"errors"
	"fmt"
	"time"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"

	"github.com/vadimfilimonov/house/internal/models"
	"github.com/vadimfilimonov/house/internal/storage/pg"
)

var (
	ErrUserNotFound = errors.New("user is not found")
	defaultTimeout  = 5 * time.Second
)

type Store struct {
	storage *pg.Storage
}

func New(storage *pg.Storage) *Store {
	return &Store{storage: storage}
}

func (s *Store) Add(ctx context.Context, email, hashedPassword, userType string) (*string, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	id := uuid.New().String()

	query := `INSERT INTO users (id, email, password, user_type) VALUES ($1, $2, $3, $4)`
	if _, err := s.storage.ExecContext(ctx, query, id, email, hashedPassword, userType); err != nil {
		return nil, fmt.Errorf("cannot add user to database: %w", err)
	}

	return &id, nil
}

func (s *Store) Get(ctx context.Context, id string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	query := "SELECT email, password, user_type FROM users WHERE id = $1 LIMIT 1"

	sqlRow := s.storage.QueryRowContext(ctx, query, id)
	if sqlRow == nil {
		return nil, fmt.Errorf("sql row is nil")
	}

	var email, hashedPassword, userType string
	if err := sqlRow.Scan(&email, &hashedPassword, &userType); err != nil {
		return nil, ErrUserNotFound
	}

	return &models.User{
		ID:       id,
		Email:    email,
		Password: hashedPassword,
		UserType: userType,
	}, nil
}

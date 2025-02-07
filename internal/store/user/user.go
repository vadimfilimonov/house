package storage

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
)

type Store struct {
	storage *pg.Storage
}

func New(storage *pg.Storage) *Store {
	return &Store{storage: storage}
}

func (d *Store) Add(ctx context.Context, email, hashedPassword, userType string) (*string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	id := uuid.New().String()

	query := `INSERT INTO users (user_id, email, password, user_type) VALUES ($1, $2, $3, $4)`
	_, err := d.storage.ExecContext(ctx, query, id, email, hashedPassword, userType)
	if err != nil {
		return nil, fmt.Errorf("cannot add user to database: %w", err)
	}

	return &id, nil
}

func (d *Store) Get(ctx context.Context, id string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var email string
	var hashedPassword string
	var userType string

	query := "SELECT email, password, user_type FROM users WHERE user_id = $1 LIMIT 1"

	sqlRow := d.storage.QueryRowContext(ctx, query, id)
	if sqlRow == nil {
		return nil, fmt.Errorf("sql row is nil")
	}

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

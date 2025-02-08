package house

import (
	"context"
	"fmt"
	"time"

	"github.com/vadimfilimonov/house/internal/models"
	"github.com/vadimfilimonov/house/internal/storage/pg"
)

var (
	defaultTimeout = 5 * time.Second
)

type Store struct {
	storage *pg.Storage
}

func New(storage *pg.Storage) *Store {
	return &Store{storage: storage}
}

func (s *Store) Add(ctx context.Context, address string, year int, developer *string) (*models.House, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	query := `INSERT INTO houses (address, year, developer, created_at) VALUES ($1, $2, $3, NOW())`

	_, err := s.storage.ExecContext(ctx, query, address, year, developer)
	if err != nil {
		return nil, fmt.Errorf("cannot add house to database: %w", err)
	}

	// TODO: Получать HouseID, CreatedAt
	return &models.House{
		Address:   address,
		Year:      year,
		Developer: developer,
	}, nil
}

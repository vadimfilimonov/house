package house

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/vadimfilimonov/house/internal/models"
	"github.com/vadimfilimonov/house/internal/storage/pg"
)

var (
	ErrHouseNotAdded = errors.New(("house is not added"))
	defaultTimeout   = 5 * time.Second
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

	query := `INSERT INTO houses (address, year, developer, created_at) VALUES ($1, $2, $3, NOW()) RETURNING id, created_at`

	sqlRow := s.storage.QueryRowContext(ctx, query, address, year, developer)
	if sqlRow == nil {
		return nil, fmt.Errorf("sql row is nil")
	}

	var houseID int
	var timestamp time.Time
	if err := sqlRow.Scan(&houseID, &timestamp); err != nil {
		return nil, ErrHouseNotAdded
	}

	createdAt := timestamp.Format("2006-01-02T15:04:05Z")

	return &models.House{
		ID:        models.HouseID(houseID),
		Address:   address,
		Year:      year,
		Developer: developer,
		CreatedAt: &createdAt,
	}, nil
}

func (s *Store) Update(ctx context.Context, houseID int) error {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	query := `UPDATE houses SET update_at = NOW() WHERE id = $1`
	if _, err := s.storage.ExecContext(ctx, query, houseID); err != nil {
		return fmt.Errorf("cannot update houses table: %w", err)
	}

	return nil
}

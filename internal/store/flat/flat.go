package flat

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

func (s *Store) Add(ctx context.Context, number, houseID, price, rooms int) (*models.Flat, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	query := `INSERT INTO flats (number, house_id, price, rooms, status) VALUES ($1, $2, $3, $4, $5)`
	if _, err := s.storage.ExecContext(ctx, query, number, houseID, price, rooms, models.CreatedStatus); err != nil {
		return nil, fmt.Errorf("cannot add flat to database: %w", err)
	}

	return &models.Flat{
		Number:  number,
		HouseID: models.HouseID(houseID),
		Price:   price,
		Rooms:   rooms,
		Status:  models.CreatedStatus,
	}, nil
}

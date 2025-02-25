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

	tx, err := s.storage.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("cannot begin transaction: %w", err)
	}

	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				err = fmt.Errorf("transaction rollback failed: %w", rollbackErr)
			}
		}
	}()

	flatsQuery := `INSERT INTO flats (number, house_id, price, rooms, status) VALUES ($1, $2, $3, $4, $5)`
	_, err = tx.ExecContext(ctx, flatsQuery, number, houseID, price, rooms, models.CreatedStatus)
	if err != nil {
		return nil, fmt.Errorf("cannot add flat to database: %w", err)
	}

	housesQuery := `UPDATE houses SET update_at = NOW() WHERE id = $1`
	_, err = tx.ExecContext(ctx, housesQuery, houseID)
	if err != nil {
		return nil, fmt.Errorf("cannot update houses table: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("transaction commit failed: %w", err)
	}

	return &models.Flat{
		Number:  number,
		HouseID: models.HouseID(houseID),
		Price:   price,
		Rooms:   rooms,
		Status:  models.CreatedStatus,
	}, nil
}

package subscription

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/lib/pq"

	"github.com/vadimfilimonov/house/internal/storage/pg"
	"github.com/vadimfilimonov/house/internal/store/pgerr"
)

var (
	ErrHouseNotFound             = errors.New("house is not found")
	ErrSubscriptionAlreadyExists = errors.New("subscription already exists")
	defaultTimeout               = 5 * time.Second
)

type Store struct {
	storage *pg.Storage
}

func New(storage *pg.Storage) *Store {
	return &Store{storage: storage}
}

func (s *Store) Add(ctx context.Context, houseID int, email string) error {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	query := `INSERT INTO house_subscriptions (house_id, email) VALUES ($1, $2)`
	if _, err := s.storage.ExecContext(ctx, query, houseID, email); err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			switch pqErr.Code {
			case pgerr.ForeignKeyViolation:
				return ErrHouseNotFound
			case pgerr.UniqueViolation:
				return ErrSubscriptionAlreadyExists
			}
		}

		return fmt.Errorf("cannot add house subscription: %w", err)
	}

	return nil
}

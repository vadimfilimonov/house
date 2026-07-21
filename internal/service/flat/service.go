package flat

import (
	"context"
	"fmt"

	"github.com/vadimfilimonov/house/internal/models"
)

const (
	minNumber     = 1
	minHouseID    = 1
	minPrice      = 0
	minRoomsCount = 1
)

type flatStore interface {
	Add(ctx context.Context, number, houseID, price, rooms int) (*models.Flat, error)
	ListByHouseID(ctx context.Context, houseID int, includeAllStatuses bool) ([]models.Flat, error)
	UpdateStatus(ctx context.Context, flatID int, status models.Status) (*models.Flat, error)
}

type Flat struct {
	store flatStore
}

func New(store flatStore) *Flat {
	return &Flat{store: store}
}

func (f *Flat) Create(ctx context.Context, number, houseID, price, rooms int) (*models.Flat, error) {
	err := validate(number, houseID, price, rooms)
	if err != nil {
		return nil, err
	}

	flat, err := f.store.Add(ctx, number, houseID, price, rooms)
	if err != nil {
		return nil, err
	}

	return flat, nil
}

func (f *Flat) ListByHouseID(ctx context.Context, houseID int, includeAllStatuses bool) ([]models.Flat, error) {
	if houseID < minHouseID {
		return nil, fmt.Errorf("houseID \"%d\" cannot be less than %d", houseID, minHouseID)
	}

	return f.store.ListByHouseID(ctx, houseID, includeAllStatuses)
}

func (f *Flat) UpdateStatus(ctx context.Context, flatID int, status models.Status) (*models.Flat, error) {
	if flatID < 1 {
		return nil, fmt.Errorf("flatID \"%d\" cannot be less than 1", flatID)
	}

	if status != models.OnModerationStatus && status != models.ApprovedStatus && status != models.DeclinedStatus {
		return nil, fmt.Errorf("status %q is not allowed for moderation update", status)
	}

	return f.store.UpdateStatus(ctx, flatID, status)
}

func validate(number, houseID, price, rooms int) error {
	if number < minNumber {
		return fmt.Errorf("flat number \"%d\" cannot be less than %d", number, minNumber)
	}

	if houseID < minHouseID {
		return fmt.Errorf("houseID \"%d\" cannot be less than %d", houseID, minHouseID)
	}

	if price < minPrice {
		return fmt.Errorf("price \"%d\" cannot be less than %d", price, minPrice)
	}

	if rooms < minRoomsCount {
		return fmt.Errorf("rooms count \"%d\" cannot be less than %d", rooms, minRoomsCount)
	}

	return nil
}

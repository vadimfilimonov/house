package flat

import (
	"context"
	"fmt"

	"github.com/vadimfilimonov/house/internal/models"
)

const (
	minNumber     = 1
	minHouseID    = 0
	minPrice      = 0
	minRoomsCount = 1
)

type flatStore interface {
	Add(ctx context.Context, number, houseID, price, rooms int) (*models.Flat, error)
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

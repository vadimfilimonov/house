package house

import (
	"context"
	"fmt"

	"github.com/vadimfilimonov/house/internal/models"
)

type houseStore interface {
	Add(ctx context.Context, address string, year int, developer *string) (*models.House, error)
	Update(ctx context.Context, houseID int) error
}

type House struct {
	store houseStore
}

func New(houseStore houseStore) *House {
	return &House{
		store: houseStore,
	}
}

func (h *House) Create(ctx context.Context, address string, year int, developer *string) (*models.House, error) {
	if err := validateYear(year); err != nil {
		return nil, err
	}

	house, err := h.store.Add(ctx, address, year, developer)
	if err != nil {
		return nil, err
	}

	return house, nil
}

func (h *House) Update(ctx context.Context, id int) error {
	return h.store.Update(ctx, id)
}

func validateYear(year int) error {
	if year < 0 {
		return fmt.Errorf("year cannot be less than 0: %d", year)
	}

	return nil
}

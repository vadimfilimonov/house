package houseget

import (
	"context"

	"github.com/vadimfilimonov/house/internal/models"
)

type flatManager interface {
	// ListByHouseID returns flats linked to the house with optional status filtering.
	ListByHouseID(ctx context.Context, houseID int, includeAllStatuses bool) ([]models.Flat, error)
}

package flatcreate

import (
	"context"

	"github.com/vadimfilimonov/house/internal/models"
)

type flatManager interface {
	// Create stores a new flat in created status.
	Create(ctx context.Context, number, houseID, price, rooms int) (*models.Flat, error)
}

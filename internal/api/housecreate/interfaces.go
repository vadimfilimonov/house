package housecreate

import (
	"context"

	"github.com/vadimfilimonov/house/internal/models"
)

type houseManager interface {
	// Create stores a new house and returns its public data.
	Create(ctx context.Context, address string, year int, developer *string) (*models.House, error)
}

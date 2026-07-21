package flatupdate

import (
	"context"

	"github.com/vadimfilimonov/house/internal/models"
)

type flatManager interface {
	// UpdateStatus changes a flat moderation status using allowed transitions.
	UpdateStatus(ctx context.Context, flatID int, status models.Status) (*models.Flat, error)
}

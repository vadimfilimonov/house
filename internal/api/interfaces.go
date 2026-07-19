package api

import (
	"context"

	"github.com/vadimfilimonov/house/internal/models"
)

type UserManager interface {
	// Register creates a user with the requested role.
	Register(ctx context.Context, email, password, userType string) (userID *string, err error)
	// Login authenticates a user and returns an authorization token.
	Login(ctx context.Context, id string, password string) (token *string, err error)
}

type HouseManager interface {
	// Create stores a new house and returns its public data.
	Create(ctx context.Context, address string, year int, developer *string) (*models.House, error)
}

type FlatManager interface {
	// Create stores a new flat in created status.
	Create(ctx context.Context, number, houseID, price, rooms int) (*models.Flat, error)
	// ListByHouseID returns flats linked to the house with optional status filtering.
	ListByHouseID(ctx context.Context, houseID int, includeAllStatuses bool) ([]models.Flat, error)
	// UpdateStatus changes a flat moderation status using allowed transitions.
	UpdateStatus(ctx context.Context, flatID int, status models.Status) (*models.Flat, error)
}

type SubscriptionManager interface {
	// Create subscribes an email to house updates.
	Create(ctx context.Context, houseID int, email string) error
}

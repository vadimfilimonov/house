package housesubscribe

import "context"

type subscriptionManager interface {
	// Create subscribes an email to house updates.
	Create(ctx context.Context, houseID int, email string) error
}

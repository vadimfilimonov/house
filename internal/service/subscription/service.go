package subscription

import (
	"context"
	"fmt"
)

type subscriptionStore interface {
	Add(ctx context.Context, houseID int, email string) error
}

type Subscription struct {
	store subscriptionStore
}

func New(store subscriptionStore) *Subscription {
	return &Subscription{store: store}
}

func (s *Subscription) Create(ctx context.Context, houseID int, email string) error {
	if houseID < 1 {
		return fmt.Errorf("houseID \"%d\" cannot be less than 1", houseID)
	}

	if email == "" {
		return fmt.Errorf("email cannot be empty")
	}

	// TODO: Add asynchronous email notifications for new flats in the subscribed house.
	return s.store.Add(ctx, houseID, email)
}

package housesubscribe

import (
	"fmt"

	"github.com/vadimfilimonov/house/internal/api"
)

type Input struct {
	// House identifier from the path.
	HouseID int
	// Email that receives house updates.
	Email string `json:"email"`
}

// Validate checks that the subscription request matches API constraints.
func (i Input) Validate() error {
	if i.HouseID < 1 {
		return fmt.Errorf("house id cannot be less than 1")
	}

	if err := api.ValidateEmail(i.Email); err != nil {
		return err
	}

	return nil
}

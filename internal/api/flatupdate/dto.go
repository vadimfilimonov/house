package flatupdate

import (
	"fmt"

	"github.com/vadimfilimonov/house/internal/models"
)

type Input struct {
	// Flat identifier.
	ID int `json:"id"`
	// New moderation status.
	Status string `json:"status"`
}

// Validate checks that the update-flat request matches moderation rules.
func (i Input) Validate() error {
	if i.ID < 1 {
		return fmt.Errorf("flat id cannot be less than 1")
	}

	status := models.Status(i.Status)
	if status != models.OnModerationStatus && status != models.ApprovedStatus && status != models.DeclinedStatus {
		return fmt.Errorf("status %q is not allowed for moderation update", i.Status)
	}

	return nil
}

package housecreate

import (
	"fmt"

	"github.com/vadimfilimonov/house/internal/api"
)

type Input struct {
	// House address.
	Address string `json:"address"`
	// House construction year.
	Year int `json:"year"`
	// House developer.
	Developer *string `json:"developer,omitempty"`
}

// Validate checks that the create-house request matches API constraints.
func (i Input) Validate() error {
	if i.Address == "" {
		return fmt.Errorf("address cannot be empty")
	}

	if err := api.ValidateMaxStringSize("address", i.Address); err != nil {
		return err
	}

	if i.Developer != nil {
		if err := api.ValidateMaxStringSize("developer", *i.Developer); err != nil {
			return err
		}
	}

	if i.Year < 0 {
		return fmt.Errorf("year cannot be less than 0")
	}

	return nil
}

type Output struct {
	// Created house identifier.
	ID int `json:"id"`
	// House address.
	Address string `json:"address"`
	// House construction year.
	Year int `json:"year"`
	// House developer.
	Developer *string `json:"developer,omitempty"`
	// House creation date.
	CreatedAt *string `json:"created_at,omitempty"`
	// Date when a flat was last added to the house.
	UpdateAt *string `json:"update_at,omitempty"`
}

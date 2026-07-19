package houseget

import (
	"fmt"

	"github.com/vadimfilimonov/house/internal/api"
)

type Input struct {
	// House identifier from the path.
	ID int
}

// Validate checks that the house-list request matches API constraints.
func (i Input) Validate() error {
	if i.ID < 1 {
		return fmt.Errorf("house id cannot be less than 1")
	}

	return nil
}

type Output struct {
	// Flats linked to the requested house.
	Flats []api.FlatOutput `json:"flats"`
}

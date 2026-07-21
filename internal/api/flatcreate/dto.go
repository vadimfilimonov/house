package flatcreate

import "fmt"

type Input struct {
	// Flat number inside the house.
	Number int `json:"number"`
	// House identifier linked to the flat.
	HouseID int `json:"house_id"`
	// Flat price in conventional units.
	Price int `json:"price"`
	// Number of rooms in the flat.
	Rooms int `json:"rooms"`
}

// Validate checks that the create-flat request matches API constraints.
func (i Input) Validate() error {
	if i.Number < 1 {
		return fmt.Errorf("number cannot be less than 1")
	}

	if i.HouseID < 1 {
		return fmt.Errorf("house_id cannot be less than 1")
	}

	if i.Price < 0 {
		return fmt.Errorf("price cannot be less than 0")
	}

	if i.Rooms < 1 {
		return fmt.Errorf("rooms cannot be less than 1")
	}

	return nil
}

type FlatOutput struct {
	// Flat identifier.
	ID int `json:"id"`
	// Flat number inside the house.
	Number int `json:"number"`
	// House identifier linked to the flat.
	HouseID int `json:"house_id"`
	// Flat price in conventional units.
	Price int `json:"price"`
	// Number of rooms in the flat.
	Rooms int `json:"rooms"`
	// Flat moderation status.
	Status string `json:"status"`
}

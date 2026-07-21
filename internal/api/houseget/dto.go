package houseget

import "fmt"

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
	Flats []FlatOutput `json:"flats"`
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

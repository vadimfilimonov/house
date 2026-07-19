package api

import "github.com/vadimfilimonov/house/internal/models"

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

// NewFlatOutput converts a flat model into the API response DTO.
func NewFlatOutput(flat models.Flat) FlatOutput {
	return FlatOutput{
		ID:      flat.ID,
		Number:  flat.Number,
		HouseID: flat.HouseID.Int(),
		Price:   flat.Price,
		Rooms:   flat.Rooms,
		Status:  flat.Status.String(),
	}
}

package api

import "github.com/vadimfilimonov/house/internal/models"

const (
	ContextKeyUser = "user"
)

type FlatOutput struct {
	ID      int    `json:"id"`       // Flat identifier.
	Number  int    `json:"number"`   // Flat number inside the house.
	HouseID int    `json:"house_id"` // House identifier linked to the flat.
	Price   int    `json:"price"`    // Flat price in conventional units.
	Rooms   int    `json:"rooms"`    // Number of rooms in the flat.
	Status  string `json:"status"`   // Flat moderation status.
}

func newFlatOutput(flat models.Flat) FlatOutput {
	return FlatOutput{
		ID:      flat.ID,
		Number:  flat.Number,
		HouseID: flat.HouseID.Int(),
		Price:   flat.Price,
		Rooms:   flat.Rooms,
		Status:  flat.Status.String(),
	}
}

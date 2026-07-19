package houseget

import "github.com/vadimfilimonov/house/internal/models"

func convToResponse(flat models.Flat) FlatOutput {
	return FlatOutput{
		ID:      flat.ID,
		Number:  flat.Number,
		HouseID: flat.HouseID.Int(),
		Price:   flat.Price,
		Rooms:   flat.Rooms,
		Status:  flat.Status.String(),
	}
}

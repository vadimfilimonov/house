package housecreate

import "github.com/vadimfilimonov/house/internal/models"

func convToResponse(house models.House) Output {
	return Output{
		ID:        house.ID.Int(),
		Address:   house.Address,
		Year:      house.Year,
		Developer: house.Developer,
		CreatedAt: house.CreatedAt,
		UpdateAt:  house.UpdateAt,
	}
}

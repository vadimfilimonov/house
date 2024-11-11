package house

import "github.com/vadimfilimonov/house/internal/models"

type House struct {
}

func New() *House {
	return &House{}
}

func (h *House) Create(address string, year int, developer *string) (*models.House, error) {
	// TODO: Add implementation
	return nil, nil
}

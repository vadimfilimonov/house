package api

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/vadimfilimonov/house/internal/models"
)

type flatManager interface {
	Create(ctx context.Context, id, houseID, price, rooms int) (*models.Flat, error)
}

type FlatCreateInput struct {
	ID      int `json:"id"`
	HouseID int `json:"house_id"`
	Price   int `json:"price"`
	Rooms   int `json:"rooms"`
}

type FlatCreateOutput struct {
	ID      int    `json:"id"`
	HouseID int    `json:"house_id"`
	Price   int    `json:"price"`
	Rooms   int    `json:"rooms"`
	Status  string `json:"status"`
}

type FlatCreate struct {
	flatManager flatManager
}

func NewFlatCreate(flatManager flatManager) *FlatCreate {
	return &FlatCreate{
		flatManager: flatManager,
	}
}

func (f *FlatCreate) Handle(c *fiber.Ctx) error {
	ctx := c.UserContext()

	_, err := jwtPayloadFromRequest(c)
	if err != nil {
		c.SendStatus(fiber.StatusUnauthorized)
		return err
	}

	var requestBody FlatCreateInput
	if err := c.BodyParser(&requestBody); err != nil {
		return fmt.Errorf("body parser: %w", err)
	}

	flat, err := f.flatManager.Create(ctx, requestBody.ID, requestBody.HouseID, requestBody.Price, requestBody.Rooms)
	if err != nil {
		c.SendStatus(fiber.StatusInternalServerError)
		return err
	}

	return c.JSON(FlatCreateOutput{
		ID:      flat.ID,
		HouseID: flat.HouseID.Int(),
		Price:   flat.Price,
		Rooms:   flat.Rooms,
		Status:  flat.Status.String(),
	})
}

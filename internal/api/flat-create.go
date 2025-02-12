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
	Number  int `json:"number"`
	HouseID int `json:"house_id"`
	Price   int `json:"price"`
	Rooms   int `json:"rooms"`
}

type FlatCreateOutput struct {
	Number  int    `json:"number"`
	HouseID int    `json:"house_id"`
	Price   int    `json:"price"`
	Rooms   int    `json:"rooms"`
	Status  string `json:"status"`
}

type FlatCreate struct {
	flatManager  flatManager
	houseManager houseManager
}

func NewFlatCreate(flatManager flatManager, houseManager houseManager) *FlatCreate {
	return &FlatCreate{
		flatManager:  flatManager,
		houseManager: houseManager,
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

	flat, err := f.flatManager.Create(ctx, requestBody.Number, requestBody.HouseID, requestBody.Price, requestBody.Rooms)
	if err != nil {
		c.SendStatus(fiber.StatusInternalServerError)
		return err
	}

	// TODO: Переделать на транзакцию
	err = f.houseManager.Update(ctx, requestBody.HouseID)
	if err != nil {
		c.SendStatus(fiber.StatusInternalServerError)
		return err
	}

	return c.JSON(FlatCreateOutput{
		Number:  flat.Number,
		HouseID: flat.HouseID.Int(),
		Price:   flat.Price,
		Rooms:   flat.Rooms,
		Status:  flat.Status.String(),
	})
}

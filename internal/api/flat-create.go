package api

import (
	"context"
	"fmt"
	"log"

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

func (i FlatCreateInput) Validate() error {
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

type FlatCreateOutput struct {
	ID      int    `json:"id"`
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
		if sendErr := c.SendStatus(fiber.StatusUnauthorized); sendErr != nil {
			log.Printf("cannot send status %d: %v", fiber.StatusUnauthorized, sendErr)
		}

		return err
	}

	var requestBody FlatCreateInput
	if err := c.BodyParser(&requestBody); err != nil {
		return fmt.Errorf("body parser: %w", err)
	}

	if err := requestBody.Validate(); err != nil {
		if sendErr := c.SendStatus(fiber.StatusBadRequest); sendErr != nil {
			log.Printf("cannot send status %d: %v", fiber.StatusBadRequest, sendErr)
		}

		return err
	}

	flat, err := f.flatManager.Create(ctx, requestBody.Number, requestBody.HouseID, requestBody.Price, requestBody.Rooms)
	if err != nil {
		if sendErr := c.SendStatus(fiber.StatusInternalServerError); sendErr != nil {
			log.Printf("cannot send status %d: %v", fiber.StatusInternalServerError, sendErr)
		}

		return err
	}

	return c.JSON(FlatCreateOutput{
		ID:      flat.ID,
		Number:  flat.Number,
		HouseID: flat.HouseID.Int(),
		Price:   flat.Price,
		Rooms:   flat.Rooms,
		Status:  flat.Status.String(),
	})
}

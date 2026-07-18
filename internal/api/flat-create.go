package api

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/vadimfilimonov/house/internal/models"
	flatStore "github.com/vadimfilimonov/house/internal/store/flat"
)

type flatManager interface {
	Create(ctx context.Context, id, houseID, price, rooms int) (*models.Flat, error)
}

type FlatCreateInput struct {
	Number  int `json:"number"`   // Номер квартиры в доме
	HouseID int `json:"house_id"` // Идентификатор дома, к которому относится квартира
	Price   int `json:"price"`    // Цена квартиры в у.е.
	Rooms   int `json:"rooms"`    // Количество комнат в квартире
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
	ID      int    `json:"id"`       // Идентификатор созданной квартиры
	Number  int    `json:"number"`   // Номер квартиры в доме
	HouseID int    `json:"house_id"` // Идентификатор дома, к которому относится квартира
	Price   int    `json:"price"`    // Цена квартиры в у.е.
	Rooms   int    `json:"rooms"`    // Количество комнат в квартире
	Status  string `json:"status"`   // Статус модерации квартиры
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
		if errors.Is(err, flatStore.ErrFlatAlreadyExists) {
			if sendErr := c.SendStatus(fiber.StatusConflict); sendErr != nil {
				log.Printf("cannot send status %d: %v", fiber.StatusConflict, sendErr)
			}

			return err
		}

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

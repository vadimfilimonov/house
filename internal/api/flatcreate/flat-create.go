package flatcreate

import (
	"errors"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/vadimfilimonov/house/internal/api"
	flatStore "github.com/vadimfilimonov/house/internal/store/flat"
)

type FlatCreate struct {
	flatManager  api.FlatManager
	houseManager api.HouseManager
}

func New(flatManager api.FlatManager, houseManager api.HouseManager) *FlatCreate {
	return &FlatCreate{
		flatManager:  flatManager,
		houseManager: houseManager,
	}
}

func (f *FlatCreate) Handle(c *fiber.Ctx) error {
	ctx := c.UserContext()

	_, err := api.JWTPayloadFromRequest(c)
	if err != nil {
		if sendErr := c.SendStatus(fiber.StatusUnauthorized); sendErr != nil {
			log.Printf("cannot send status %d: %v", fiber.StatusUnauthorized, sendErr)
		}

		return err
	}

	var requestBody Input
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

	return c.JSON(api.NewFlatOutput(*flat))
}

package houseget

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/vadimfilimonov/house/internal/api"
	"github.com/vadimfilimonov/house/internal/models"
	"github.com/vadimfilimonov/house/internal/service/auth_token"
)

type flatManager interface {
	// ListByHouseID returns flats linked to the house with optional status filtering.
	ListByHouseID(ctx context.Context, houseID int, includeAllStatuses bool) ([]models.Flat, error)
}

type HouseGet struct {
	flatManager flatManager
}

func New(flatManager flatManager) *HouseGet {
	return &HouseGet{flatManager: flatManager}
}

func (h *HouseGet) Handle(c *fiber.Ctx) error {
	ctx := c.UserContext()

	jwtPayload, err := api.JWTPayloadFromRequest(c)
	if err != nil {
		if sendErr := c.SendStatus(fiber.StatusUnauthorized); sendErr != nil {
			log.Printf("cannot send status %d: %v", fiber.StatusUnauthorized, sendErr)
		}

		return err
	}

	userType, ok := jwtPayload[auth_token.ClaimsKeyUserType].(string)
	if !ok {
		err := fmt.Errorf("cannot get user type")
		if sendErr := c.SendStatus(fiber.StatusInternalServerError); sendErr != nil {
			log.Printf("cannot send status %d: %v", fiber.StatusInternalServerError, sendErr)
		}

		return err
	}

	if userType != models.UserTypeClient && userType != models.UserTypeModerator {
		err := fmt.Errorf("user type %q cannot get house flats", userType)
		if sendErr := c.SendStatus(fiber.StatusForbidden); sendErr != nil {
			log.Printf("cannot send status %d: %v", fiber.StatusForbidden, sendErr)
		}

		return err
	}

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		if sendErr := c.SendStatus(fiber.StatusBadRequest); sendErr != nil {
			log.Printf("cannot send status %d: %v", fiber.StatusBadRequest, sendErr)
		}

		return fmt.Errorf("house id is not valid: %w", err)
	}

	request := Input{ID: id}
	if err := request.Validate(); err != nil {
		if sendErr := c.SendStatus(fiber.StatusBadRequest); sendErr != nil {
			log.Printf("cannot send status %d: %v", fiber.StatusBadRequest, sendErr)
		}

		return err
	}

	includeAllStatuses := userType == models.UserTypeModerator
	flats, err := h.flatManager.ListByHouseID(ctx, request.ID, includeAllStatuses)
	if err != nil {
		if sendErr := c.SendStatus(fiber.StatusInternalServerError); sendErr != nil {
			log.Printf("cannot send status %d: %v", fiber.StatusInternalServerError, sendErr)
		}

		return err
	}

	output := Output{Flats: make([]api.FlatOutput, 0, len(flats))}
	for _, flat := range flats {
		output.Flats = append(output.Flats, api.NewFlatOutput(flat))
	}

	return c.JSON(output)
}

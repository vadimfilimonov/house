package api

import (
	"context"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/vadimfilimonov/house/internal/models"
	"github.com/vadimfilimonov/house/internal/service/auth_token"
)

type houseManager interface {
	Create(ctx context.Context, address string, year int, developer *string) (*models.House, error)
}

type HouseCreateInput struct {
	Address   string  `json:"address"`             // House address.
	Year      int     `json:"year"`                // House construction year.
	Developer *string `json:"developer,omitempty"` // House developer.
}

// Validate checks that the create-house request matches API constraints.
func (i HouseCreateInput) Validate() error {
	if i.Address == "" {
		return fmt.Errorf("address cannot be empty")
	}

	if err := validateMaxStringSize("address", i.Address); err != nil {
		return err
	}

	if i.Developer != nil {
		if err := validateMaxStringSize("developer", *i.Developer); err != nil {
			return err
		}
	}

	if i.Year < 0 {
		return fmt.Errorf("year cannot be less than 0")
	}

	return nil
}

type HouseCreateOutput struct {
	ID        int     `json:"id"`                   // Created house identifier.
	Address   string  `json:"address"`              // House address.
	Year      int     `json:"year"`                 // House construction year.
	Developer *string `json:"developer,omitempty"`  // House developer.
	CreatedAt *string `json:"created_at,omitempty"` // House creation date.
	UpdateAt  *string `json:"update_at,omitempty"`  // Date when a flat was last added to the house.
}

type HouseCreate struct {
	houseManager houseManager
}

func NewHouseCreate(houseManager houseManager) *HouseCreate {
	return &HouseCreate{
		houseManager: houseManager,
	}
}

func (h *HouseCreate) Handle(c *fiber.Ctx) error {
	ctx := c.UserContext()

	jwtPayload, err := jwtPayloadFromRequest(c)
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

	if userType != models.UserTypeModerator {
		err := fmt.Errorf("user type \"%s\" cannot create house", userType)
		if sendErr := c.SendStatus(fiber.StatusForbidden); sendErr != nil {
			log.Printf("cannot send status %d: %v", fiber.StatusForbidden, sendErr)
		}

		return err
	}

	var requestBody HouseCreateInput
	if err := c.BodyParser(&requestBody); err != nil {
		return fmt.Errorf("body parser: %w", err)
	}

	if err := requestBody.Validate(); err != nil {
		if sendErr := c.SendStatus(fiber.StatusBadRequest); sendErr != nil {
			log.Printf("cannot send status %d: %v", fiber.StatusBadRequest, sendErr)
		}

		return err
	}

	house, err := h.houseManager.Create(ctx, requestBody.Address, requestBody.Year, requestBody.Developer)
	if err != nil {
		if sendErr := c.SendStatus(fiber.StatusInternalServerError); sendErr != nil {
			log.Printf("cannot send status %d: %v", fiber.StatusInternalServerError, sendErr)
		}

		return err
	}

	return c.JSON(HouseCreateOutput{
		ID:        house.ID.Int(),
		Address:   house.Address,
		Year:      house.Year,
		Developer: house.Developer,
		CreatedAt: house.CreatedAt,
		UpdateAt:  house.UpdateAt,
	})
}

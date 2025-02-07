package api

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/vadimfilimonov/house/internal/models"
	"github.com/vadimfilimonov/house/internal/service/auth_token"
)

type houseManager interface {
	Create(address string, year int, developer *string) (*models.House, error)
}

type HouseCreateInput struct {
	Address   string  `json:"address"`
	Year      int     `json:"year"`
	Developer *string `json:"developer,omitempty"`
}

type HouseCreateOutput struct {
	ID        int     `json:"id"`
	Address   string  `json:"address"`
	Year      int     `json:"year"`
	Developer *string `json:"developer,omitempty"`
	CreatedAt *string `json:"created_at,omitempty"`
	UpdateAt  *string `json:"update_at,omitempty"`
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
	jwtPayload, err := jwtPayloadFromRequest(c)
	if err != nil {
		c.SendStatus(fiber.StatusUnauthorized)
		return err
	}

	userType, ok := jwtPayload[auth_token.ClaimsKeyUserType].(string)
	if !ok {
		c.SendStatus(fiber.StatusInternalServerError)
		return err
	}

	if userType != models.UserTypeModerator {
		c.SendStatus(fiber.StatusForbidden)
		return fmt.Errorf("user type %s cannot create house", userType)
	}

	var requestBody HouseCreateInput
	if err := c.BodyParser(&requestBody); err != nil {
		return fmt.Errorf("body parser: %w", err)
	}

	house, err := h.houseManager.Create(requestBody.Address, requestBody.Year, requestBody.Developer)
	if err != nil {
		c.SendStatus(fiber.StatusInternalServerError)
		return err
	}

	return c.JSON(HouseCreateOutput{
		ID:        house.ID,
		Address:   house.Address,
		Year:      house.Year,
		Developer: house.Developer,
		CreatedAt: house.CreatedAt,
		UpdateAt:  house.UpdateAt,
	})
}

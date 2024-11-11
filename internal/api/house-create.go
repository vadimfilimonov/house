package api

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type HouseCreateInput struct {
	Address   string  `json:"address"`
	Year      string  `json:"year"`
	Developer *string `json:"developer,omitempty"`
}

type HouseCreateOutput struct {
	ID        string  `json:"id"`
	Address   string  `json:"address"`
	Year      string  `json:"year"`
	Developer *string `json:"developer,omitempty"`
	CreatedAt *string `json:"created_at,omitempty"`
	UpdateAt  *string `json:"update_at,omitempty"`
}

type HouseCreate struct {
}

func NewHouseCreate() *HouseCreate {
	return &HouseCreate{}
}

func (h *HouseCreate) Handle(c *fiber.Ctx) error {
	var requestBody HouseCreateInput
	if err := c.BodyParser(&requestBody); err != nil {
		return fmt.Errorf("body parser: %w", err)
	}

	return c.JSON(HouseCreateOutput{})
}

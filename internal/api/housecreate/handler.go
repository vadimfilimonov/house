package housecreate

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	houseManager houseManager
}

func New(houseManager houseManager) *Handler {
	return &Handler{
		houseManager: houseManager,
	}
}

func (h *Handler) Handle(c *fiber.Ctx) error {
	ctx := c.UserContext()

	var requestBody Input
	if err := c.BodyParser(&requestBody); err != nil {
		return fmt.Errorf("body parser: %w", err)
	}

	if err := requestBody.Validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	house, err := h.houseManager.Create(ctx, requestBody.Address, requestBody.Year, requestBody.Developer)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.JSON(convToResponse(*house))
}

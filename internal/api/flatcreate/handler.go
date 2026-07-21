package flatcreate

import (
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/vadimfilimonov/house/internal/api"
	flatStore "github.com/vadimfilimonov/house/internal/store/flat"
)

type Handler struct {
	flatManager flatManager
}

func New(flatManager flatManager) *Handler {
	return &Handler{
		flatManager: flatManager,
	}
}

func (h *Handler) Handle(c *fiber.Ctx) error {
	ctx := c.UserContext()

	_, err := api.JWTPayloadFromRequest(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).SendString(err.Error())
	}

	var requestBody Input
	if err := c.BodyParser(&requestBody); err != nil {
		return fmt.Errorf("body parser: %w", err)
	}

	if err := requestBody.Validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	flat, err := h.flatManager.Create(ctx, requestBody.Number, requestBody.HouseID, requestBody.Price, requestBody.Rooms)
	if err != nil {
		if errors.Is(err, flatStore.ErrFlatAlreadyExists) {
			return c.Status(fiber.StatusConflict).SendString(err.Error())
		}

		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.JSON(convToResponse(*flat))
}

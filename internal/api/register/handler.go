package register

import (
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"

	manager "github.com/vadimfilimonov/house/internal/service/user"
)

type Handler struct {
	userManager userManager
}

func New(userManager userManager) *Handler {
	return &Handler{
		userManager: userManager,
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

	userID, err := h.userManager.Register(ctx, requestBody.Email, requestBody.Password, requestBody.UserType)
	if err != nil {
		if errors.Is(err, manager.ErrIncorrectInput) {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		return err
	}

	if userID == nil {
		return c.Status(fiber.StatusBadRequest).SendString("userID is empty")
	}

	c.Status(fiber.StatusOK)
	return c.JSON(convToResponse(*userID))
}

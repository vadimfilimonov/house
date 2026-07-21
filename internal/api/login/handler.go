package login

import (
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"

	manager "github.com/vadimfilimonov/house/internal/service/user"
	store "github.com/vadimfilimonov/house/internal/store/user"
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

	token, err := h.userManager.Login(ctx, requestBody.Email, requestBody.Password)
	if err != nil {
		if errors.Is(err, store.ErrUserNotFound) {
			return c.Status(fiber.StatusNotFound).SendString(err.Error())
		}

		if errors.Is(err, manager.ErrWrongPassword) {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		return err
	}

	c.Set("Content-Type", "application/json")

	return c.JSON(convToResponse(*token))
}

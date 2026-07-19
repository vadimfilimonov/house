package login

import (
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/vadimfilimonov/house/internal/api"

	manager "github.com/vadimfilimonov/house/internal/service/user"
	store "github.com/vadimfilimonov/house/internal/store/user"
)

type Login struct {
	userManager api.UserManager
}

func New(userManager api.UserManager) *Login {
	return &Login{
		userManager: userManager,
	}
}

func (h *Login) Handle(c *fiber.Ctx) error {
	ctx := c.UserContext()

	var requestBody Input
	if err := c.BodyParser(&requestBody); err != nil {
		return fmt.Errorf("body parser: %w", err)
	}

	if err := requestBody.Validate(); err != nil {
		c.Status(fiber.StatusBadRequest)
		return err
	}

	token, err := h.userManager.Login(ctx, requestBody.Email, requestBody.Password)
	if err != nil {
		if errors.Is(err, store.ErrUserNotFound) {
			c.Status(fiber.StatusNotFound)
			return err
		}

		if errors.Is(err, manager.ErrWrongPassword) {
			c.Status(fiber.StatusBadRequest)
			return err
		}

		return err
	}

	c.Set("Content-Type", "application/json")

	return c.JSON(Output{Token: *token})
}

package api

import (
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"

	manager "github.com/vadimfilimonov/house/internal/service/user"
	store "github.com/vadimfilimonov/house/internal/store/user"
)

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginOutput struct {
	Token string `json:"token"`
}

type Login struct {
	userManager userManager
}

func NewLogin(userManager userManager) *Login {
	return &Login{
		userManager: userManager,
	}
}

func (h *Login) Handle(c *fiber.Ctx) error {
	ctx := c.UserContext()

	var requestBody LoginInput
	if err := c.BodyParser(&requestBody); err != nil {
		return fmt.Errorf("body parser: %w", err)
	}

	if requestBody.Email == "" {
		c.Status(fiber.StatusBadRequest)
		return fmt.Errorf("email cannot be empty")
	}

	if requestBody.Password == "" {
		c.Status(fiber.StatusBadRequest)
		return fmt.Errorf("password cannot be empty")
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

	return c.JSON(LoginOutput{Token: *token})
}

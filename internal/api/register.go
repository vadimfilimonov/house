package api

import (
	"context"
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
	manager "github.com/vadimfilimonov/house/internal/service/user"
)

type RegisterInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	UserType string `json:"user_type"`
}

type RegisterOutput struct {
	UserID string `json:"user_id"`
}

type Register struct {
	userManager userManager
}

func NewRegister(userManager userManager) *Register {
	return &Register{
		userManager: userManager,
	}
}

func (h *Register) Handle(c *fiber.Ctx) error {
	ctx := context.Background()

	var requestBody RegisterInput
	if err := c.BodyParser(&requestBody); err != nil {
		return fmt.Errorf("body parser: %w", err)
	}

	userID, err := h.userManager.Register(ctx, requestBody.Email, requestBody.Password, requestBody.UserType)
	if err != nil {
		if errors.Is(err, manager.ErrIncorrectInput) {
			c.SendStatus(fiber.StatusBadRequest)
			return err
		}

		return err
	}

	if userID == nil {
		c.SendStatus(fiber.StatusBadRequest)
		return fmt.Errorf("userID is empty")
	}

	c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	c.SendStatus(fiber.StatusCreated)

	return c.JSON(RegisterOutput{UserID: *userID})
}

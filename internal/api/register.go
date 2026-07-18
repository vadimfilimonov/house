package api

import (
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

func (i RegisterInput) Validate() error {
	if err := validateEmail(i.Email); err != nil {
		return err
	}

	if err := validatePassword(i.Password); err != nil {
		return err
	}

	if err := validateUserType(i.UserType); err != nil {
		return err
	}

	return nil
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
	ctx := c.UserContext()

	var requestBody RegisterInput
	if err := c.BodyParser(&requestBody); err != nil {
		return fmt.Errorf("body parser: %w", err)
	}

	if err := requestBody.Validate(); err != nil {
		c.SendStatus(fiber.StatusBadRequest)
		return err
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

	c.Status(fiber.StatusOK)
	return c.JSON(RegisterOutput{UserID: *userID})
}

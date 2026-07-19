package api

import (
	"errors"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"

	manager "github.com/vadimfilimonov/house/internal/service/user"
)

type RegisterInput struct {
	Email    string `json:"email"`     // User email.
	Password string `json:"password"`  // User password.
	UserType string `json:"user_type"` // User type: client or moderator.
}

// Validate checks that the registration request matches API constraints.
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
	UserID string `json:"user_id"` // Created user identifier.
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
		if sendErr := c.SendStatus(fiber.StatusBadRequest); sendErr != nil {
			log.Printf("cannot send status %d: %v", fiber.StatusBadRequest, sendErr)
		}

		return err
	}

	userID, err := h.userManager.Register(ctx, requestBody.Email, requestBody.Password, requestBody.UserType)
	if err != nil {
		if errors.Is(err, manager.ErrIncorrectInput) {
			if sendErr := c.SendStatus(fiber.StatusBadRequest); sendErr != nil {
				log.Printf("cannot send status %d: %v", fiber.StatusBadRequest, sendErr)
			}

			return err
		}

		return err
	}

	if userID == nil {
		if sendErr := c.SendStatus(fiber.StatusBadRequest); sendErr != nil {
			log.Printf("cannot send status %d: %v", fiber.StatusBadRequest, sendErr)
		}

		return fmt.Errorf("userID is empty")
	}

	c.Status(fiber.StatusOK)
	return c.JSON(RegisterOutput{UserID: *userID})
}

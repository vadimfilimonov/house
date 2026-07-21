package register

import (
	"errors"
	"fmt"
	"log"

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
	return c.JSON(convToResponse(*userID))
}

package housecreate

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/vadimfilimonov/house/internal/api"
	"github.com/vadimfilimonov/house/internal/models"
	"github.com/vadimfilimonov/house/internal/service/auth_token"
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

	jwtPayload, err := api.JWTPayloadFromRequest(c)
	if err != nil {
		if sendErr := c.SendStatus(fiber.StatusUnauthorized); sendErr != nil {
			log.Printf("cannot send status %d: %v", fiber.StatusUnauthorized, sendErr)
		}

		return err
	}

	userType, ok := jwtPayload[auth_token.ClaimsKeyUserType].(string)
	if !ok {
		err := fmt.Errorf("cannot get user type")
		if sendErr := c.SendStatus(fiber.StatusInternalServerError); sendErr != nil {
			log.Printf("cannot send status %d: %v", fiber.StatusInternalServerError, sendErr)
		}

		return err
	}

	if userType != models.UserTypeModerator {
		err := fmt.Errorf("user type \"%s\" cannot create house", userType)
		if sendErr := c.SendStatus(fiber.StatusForbidden); sendErr != nil {
			log.Printf("cannot send status %d: %v", fiber.StatusForbidden, sendErr)
		}

		return err
	}

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

	house, err := h.houseManager.Create(ctx, requestBody.Address, requestBody.Year, requestBody.Developer)
	if err != nil {
		if sendErr := c.SendStatus(fiber.StatusInternalServerError); sendErr != nil {
			log.Printf("cannot send status %d: %v", fiber.StatusInternalServerError, sendErr)
		}

		return err
	}

	return c.JSON(convToResponse(*house))
}

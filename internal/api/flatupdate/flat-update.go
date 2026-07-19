package flatupdate

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/vadimfilimonov/house/internal/api"
	"github.com/vadimfilimonov/house/internal/models"
	"github.com/vadimfilimonov/house/internal/service/auth_token"
	flatStore "github.com/vadimfilimonov/house/internal/store/flat"
)

type flatManager interface {
	// UpdateStatus changes a flat moderation status using allowed transitions.
	UpdateStatus(ctx context.Context, flatID int, status models.Status) (*models.Flat, error)
}

type FlatUpdate struct {
	flatManager flatManager
}

func New(flatManager flatManager) *FlatUpdate {
	return &FlatUpdate{flatManager: flatManager}
}

func (f *FlatUpdate) Handle(c *fiber.Ctx) error {
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
		err := fmt.Errorf("user type %q cannot update flat status", userType)
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

	flat, err := f.flatManager.UpdateStatus(ctx, requestBody.ID, models.Status(requestBody.Status))
	if err != nil {
		if errors.Is(err, flatStore.ErrFlatNotFound) {
			if sendErr := c.SendStatus(fiber.StatusNotFound); sendErr != nil {
				log.Printf("cannot send status %d: %v", fiber.StatusNotFound, sendErr)
			}

			return err
		}

		if errors.Is(err, flatStore.ErrFlatStatusConflict) {
			if sendErr := c.SendStatus(fiber.StatusConflict); sendErr != nil {
				log.Printf("cannot send status %d: %v", fiber.StatusConflict, sendErr)
			}

			return err
		}

		if sendErr := c.SendStatus(fiber.StatusInternalServerError); sendErr != nil {
			log.Printf("cannot send status %d: %v", fiber.StatusInternalServerError, sendErr)
		}

		return err
	}

	return c.JSON(api.NewFlatOutput(*flat))
}

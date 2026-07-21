package flatupdate

import (
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/vadimfilimonov/house/internal/api"
	"github.com/vadimfilimonov/house/internal/models"
	"github.com/vadimfilimonov/house/internal/service/auth_token"
	flatStore "github.com/vadimfilimonov/house/internal/store/flat"
)

type Handler struct {
	flatManager flatManager
}

func New(flatManager flatManager) *Handler {
	return &Handler{flatManager: flatManager}
}

func (h *Handler) Handle(c *fiber.Ctx) error {
	ctx := c.UserContext()

	jwtPayload, err := api.JWTPayloadFromRequest(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).SendString(err.Error())
	}

	userType, ok := jwtPayload[auth_token.ClaimsKeyUserType].(string)
	if !ok {
		err := fmt.Errorf("cannot get user type")
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	if userType != models.UserTypeModerator {
		err := fmt.Errorf("user type %q cannot update flat status", userType)
		return c.Status(fiber.StatusForbidden).SendString(err.Error())
	}

	var requestBody Input
	if err := c.BodyParser(&requestBody); err != nil {
		return fmt.Errorf("body parser: %w", err)
	}

	if err := requestBody.Validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	flat, err := h.flatManager.UpdateStatus(ctx, requestBody.ID, models.Status(requestBody.Status))
	if err != nil {
		if errors.Is(err, flatStore.ErrFlatNotFound) {
			return c.Status(fiber.StatusNotFound).SendString(err.Error())
		}

		if errors.Is(err, flatStore.ErrFlatStatusConflict) {
			return c.Status(fiber.StatusConflict).SendString(err.Error())
		}

		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.JSON(convToResponse(*flat))
}

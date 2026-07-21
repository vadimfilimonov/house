package housecreate

import (
	"fmt"

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
		return c.Status(fiber.StatusUnauthorized).SendString(err.Error())
	}

	userType, ok := jwtPayload[auth_token.ClaimsKeyUserType].(string)
	if !ok {
		err := fmt.Errorf("cannot get user type")
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	if userType != models.UserTypeModerator {
		err := fmt.Errorf("user type \"%s\" cannot create house", userType)
		return c.Status(fiber.StatusForbidden).SendString(err.Error())
	}

	var requestBody Input
	if err := c.BodyParser(&requestBody); err != nil {
		return fmt.Errorf("body parser: %w", err)
	}

	if err := requestBody.Validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	house, err := h.houseManager.Create(ctx, requestBody.Address, requestBody.Year, requestBody.Developer)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.JSON(convToResponse(*house))
}

package houseget

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/vadimfilimonov/house/internal/api"
	"github.com/vadimfilimonov/house/internal/models"
	"github.com/vadimfilimonov/house/internal/service/auth_token"
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

	if userType != models.UserTypeClient && userType != models.UserTypeModerator {
		err := fmt.Errorf("user type %q cannot get house flats", userType)
		return c.Status(fiber.StatusForbidden).SendString(err.Error())
	}

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(fmt.Errorf("house id is not valid: %w", err).Error())
	}

	request := Input{ID: id}
	if err := request.Validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	includeAllStatuses := userType == models.UserTypeModerator
	flats, err := h.flatManager.ListByHouseID(ctx, request.ID, includeAllStatuses)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	output := Output{Flats: make([]FlatOutput, 0, len(flats))}
	for _, flat := range flats {
		output.Flats = append(output.Flats, convToResponse(flat))
	}

	return c.JSON(output)
}

package houseget

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	authmiddleware "github.com/vadimfilimonov/house/internal/api/auth"
	"github.com/vadimfilimonov/house/internal/models"
)

type Handler struct {
	flatManager flatManager
}

func New(flatManager flatManager) *Handler {
	return &Handler{flatManager: flatManager}
}

func (h *Handler) Handle(c *fiber.Ctx) error {
	ctx := c.UserContext()

	authContext, err := authmiddleware.FromContext(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(fmt.Errorf("house id is not valid: %w", err).Error())
	}

	request := Input{ID: id}
	if err := request.Validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	includeAllStatuses := authContext.UserType == models.UserTypeModerator
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

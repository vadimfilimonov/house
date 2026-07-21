package housesubscribe

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	subscriptionStore "github.com/vadimfilimonov/house/internal/store/subscription"
)

type Handler struct {
	subscriptionManager subscriptionManager
}

func New(subscriptionManager subscriptionManager) *Handler {
	return &Handler{subscriptionManager: subscriptionManager}
}

func (h *Handler) Handle(c *fiber.Ctx) error {
	ctx := c.UserContext()

	houseID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(fmt.Errorf("house id is not valid: %w", err).Error())
	}

	var requestBody Input
	if err := c.BodyParser(&requestBody); err != nil {
		return fmt.Errorf("body parser: %w", err)
	}
	requestBody.HouseID = houseID

	if err := requestBody.Validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	if err := h.subscriptionManager.Create(ctx, requestBody.HouseID, requestBody.Email); err != nil {
		if errors.Is(err, subscriptionStore.ErrHouseNotFound) {
			return c.Status(fiber.StatusNotFound).SendString(err.Error())
		}

		if errors.Is(err, subscriptionStore.ErrSubscriptionAlreadyExists) {
			return c.Status(fiber.StatusConflict).SendString(err.Error())
		}

		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.SendStatus(fiber.StatusOK)
}

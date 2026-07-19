package housesubscribe

import (
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/vadimfilimonov/house/internal/api"
	subscriptionStore "github.com/vadimfilimonov/house/internal/store/subscription"
)

type HouseSubscribe struct {
	subscriptionManager api.SubscriptionManager
}

func New(subscriptionManager api.SubscriptionManager) *HouseSubscribe {
	return &HouseSubscribe{subscriptionManager: subscriptionManager}
}

func (h *HouseSubscribe) Handle(c *fiber.Ctx) error {
	ctx := c.UserContext()

	if _, err := api.JWTPayloadFromRequest(c); err != nil {
		if sendErr := c.SendStatus(fiber.StatusUnauthorized); sendErr != nil {
			log.Printf("cannot send status %d: %v", fiber.StatusUnauthorized, sendErr)
		}

		return err
	}

	houseID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		if sendErr := c.SendStatus(fiber.StatusBadRequest); sendErr != nil {
			log.Printf("cannot send status %d: %v", fiber.StatusBadRequest, sendErr)
		}

		return fmt.Errorf("house id is not valid: %w", err)
	}

	var requestBody Input
	if err := c.BodyParser(&requestBody); err != nil {
		return fmt.Errorf("body parser: %w", err)
	}
	requestBody.HouseID = houseID

	if err := requestBody.Validate(); err != nil {
		if sendErr := c.SendStatus(fiber.StatusBadRequest); sendErr != nil {
			log.Printf("cannot send status %d: %v", fiber.StatusBadRequest, sendErr)
		}

		return err
	}

	if err := h.subscriptionManager.Create(ctx, requestBody.HouseID, requestBody.Email); err != nil {
		if errors.Is(err, subscriptionStore.ErrHouseNotFound) {
			if sendErr := c.SendStatus(fiber.StatusNotFound); sendErr != nil {
				log.Printf("cannot send status %d: %v", fiber.StatusNotFound, sendErr)
			}

			return err
		}

		if errors.Is(err, subscriptionStore.ErrSubscriptionAlreadyExists) {
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

	return c.SendStatus(fiber.StatusOK)
}

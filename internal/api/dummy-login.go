package api

import (
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/vadimfilimonov/house/internal/models"
)

type tokenManager interface {
	Encode(sub, userType string) (*string, error)
}

type tokenStorage interface {
	Add(ctx context.Context, key, value string, expiration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
}

type DummyLoginInput struct {
	UserType string `json:"user_type"`
}

type DummyLoginOutput struct {
	Token string `json:"token"`
}

type DummyLogin struct {
	tokenManager tokenManager
	tokenStorage tokenStorage
}

func NewDummyLogin(tokenManager tokenManager, tokenStorage tokenStorage) *DummyLogin {
	return &DummyLogin{
		tokenManager: tokenManager,
		tokenStorage: tokenStorage,
	}
}

func (h *DummyLogin) Handle(c *fiber.Ctx) error {
	ctx := c.UserContext()

	var requestBody DummyLoginInput
	if err := c.BodyParser(&requestBody); err != nil {
		return fmt.Errorf("body parser: %w", err)
	}

	var userID string
	switch requestBody.UserType {
	case models.UserTypeClient:
		userID = models.FakeClientUserID
	case models.UserTypeModerator:
		userID = models.FakeModeratorUserID
	default:
		c.SendStatus(fiber.StatusNotFound)
		return fmt.Errorf("user type %s is not supported", requestBody.UserType)
	}

	var token *string

	savedToken, err := h.tokenStorage.Get(ctx, userID)
	if err == nil {
		token = &savedToken
	} else {
		token, err = h.tokenManager.Encode(userID, requestBody.UserType)
		if err != nil {
			return err
		} else if token == nil {
			return fmt.Errorf("token is nil")
		}

		if err := h.tokenStorage.Add(ctx, userID, *token, 24*time.Hour); err != nil {
			return err
		}
	}

	c.Set("Content-Type", "application/json")
	c.SendStatus(fiber.StatusOK)

	return c.JSON(DummyLoginOutput{Token: *token})
}

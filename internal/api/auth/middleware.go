package auth

import (
	"fmt"

	"github.com/gofiber/fiber/v2"

	"github.com/vadimfilimonov/house/internal/api"
	"github.com/vadimfilimonov/house/internal/models"
	"github.com/vadimfilimonov/house/internal/service/auth_token"
)

const ContextKey = "auth_context"

type Context struct {
	UserID   string
	UserType string
}

// Middleware extracts verified JWT claims once for all authorized handlers.
func Middleware(c *fiber.Ctx) error {
	payload, err := api.JWTPayloadFromRequest(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).SendString(err.Error())
	}

	userID, userIDOK := payload[auth_token.ClaimsKeySub].(string)
	userType, userTypeOK := payload[auth_token.ClaimsKeyUserType].(string)
	if !userIDOK || userID == "" || !userTypeOK || userType == "" {
		return c.Status(fiber.StatusUnauthorized).SendString("invalid authentication claims")
	}
	if userType != models.UserTypeClient && userType != models.UserTypeModerator {
		return c.Status(fiber.StatusForbidden).SendString("unsupported user role")
	}

	c.Locals(ContextKey, Context{UserID: userID, UserType: userType})
	return c.Next()
}

func FromContext(c *fiber.Ctx) (Context, error) {
	authContext, ok := c.Locals(ContextKey).(Context)
	if !ok {
		return Context{}, fmt.Errorf("authentication context is missing")
	}

	return authContext, nil
}

// RequireModerator restricts a route to users with the moderator role.
func RequireModerator(c *fiber.Ctx) error {
	authContext, err := FromContext(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	if authContext.UserType != models.UserTypeModerator {
		return c.Status(fiber.StatusForbidden).SendString("moderator role required")
	}

	return c.Next()
}

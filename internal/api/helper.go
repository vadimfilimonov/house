package api

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	jwt "github.com/golang-jwt/jwt/v5"
)

func jwtPayloadFromRequest(c *fiber.Ctx) (jwt.MapClaims, error) {
	jwtToken, ok := c.Context().Value(ContextKeyUser).(*jwt.Token)
	if !ok {
		return nil, fmt.Errorf("wrong type of JWT token in context")
	}

	payload, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("wrong type of JWT token claims")
	}

	return payload, nil
}

package main

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"

	"github.com/vadimfilimonov/house/internal/api"
	"github.com/vadimfilimonov/house/internal/service/auth_token"
	"github.com/vadimfilimonov/house/internal/service/config"
	"github.com/vadimfilimonov/house/internal/service/user"
	tokenStore "github.com/vadimfilimonov/house/internal/store/token"
	userStore "github.com/vadimfilimonov/house/internal/store/user"
)

func main() {
	ctx := context.Background()

	c := config.New()
	if err := c.Parse(); err != nil {
		log.Fatal(err)
	}

	uStorage, err := userStore.New(c.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}

	defer uStorage.Close()

	tStorage, err := tokenStore.New(c.RedisAddress, c.RedisPassword)
	if err != nil {
		log.Fatal(err)
	}

	defer tStorage.Close()

	tokenManager := auth_token.NewToken(c.JwtSecretKey)
	userManager := user.New(uStorage, tStorage, tokenManager)

	webApp := fiber.New()
	webApp.Use(contextMiddleware(ctx))
	webApp.Post("/dummyLogin", api.NewDummyLogin(tokenManager, tStorage).Handle)
	webApp.Post("/login", api.NewLogin(userManager).Handle)
	webApp.Post("/register", api.NewRegister(userManager).Handle)

	if err := webApp.Listen(c.ServerAddress); err != nil {
		log.Fatal(err)
	}
}

func contextMiddleware(ctx context.Context) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		c.SetUserContext(ctx)
		return c.Next()
	}
}

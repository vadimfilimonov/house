package main

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"

	"github.com/vadimfilimonov/house/internal/api"
	"github.com/vadimfilimonov/house/internal/service/auth_token"
	"github.com/vadimfilimonov/house/internal/service/config"
	"github.com/vadimfilimonov/house/internal/service/user"
	"github.com/vadimfilimonov/house/internal/storage/pg"
	"github.com/vadimfilimonov/house/internal/storage/redis"
	tokenStore "github.com/vadimfilimonov/house/internal/store/token"
	userStore "github.com/vadimfilimonov/house/internal/store/user"
)

func main() {
	ctx := context.Background()

	c := config.New()
	if err := c.Parse(); err != nil {
		log.Fatal(err)
	}

	database, err := pg.New(c.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	redisClient, err := redis.New(ctx, c.RedisAddress, c.RedisPassword)
	if err != nil {
		log.Fatal(err)
	}
	defer redisClient.Close()

	uStore := userStore.New(database)
	tStore := tokenStore.New(redisClient)

	tokenManager := auth_token.NewToken(c.JwtSecretKey)
	userManager := user.New(uStore, tStore, tokenManager)

	webApp := fiber.New()
	webApp.Use(contextMiddleware(ctx))
	webApp.Post("/dummyLogin", api.NewDummyLogin(tokenManager, tStore).Handle)
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

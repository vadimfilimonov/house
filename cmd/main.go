package main

import (
	"context"
	"log"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"

	"github.com/vadimfilimonov/house/internal/api"
	"github.com/vadimfilimonov/house/internal/service/auth_token"
	"github.com/vadimfilimonov/house/internal/service/config"
	"github.com/vadimfilimonov/house/internal/service/house"
	"github.com/vadimfilimonov/house/internal/service/user"
	"github.com/vadimfilimonov/house/internal/storage/pg"
	"github.com/vadimfilimonov/house/internal/storage/redis"
	houseStore "github.com/vadimfilimonov/house/internal/store/house"
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
	hStore := houseStore.New(database)
	tStore := tokenStore.New(redisClient)

	tokenManager := auth_token.NewToken(c.JwtSecretKey)
	userManager := user.New(uStore, tStore, tokenManager)
	houseManager := house.New(hStore)

	app := fiber.New()
	app.Use(contextMiddleware(ctx))

	publicGroup := app.Group("")
	publicGroup.Post("/dummyLogin", api.NewDummyLogin(tokenManager, tStore).Handle)
	publicGroup.Post("/login", api.NewLogin(userManager).Handle)
	publicGroup.Post("/register", api.NewRegister(userManager).Handle)

	authorizedGroup := app.Group("")
	authorizedGroup.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{
			Key: c.JwtSecretKey,
		},
		ContextKey: api.ContextKeyUser,
	}))
	authorizedGroup.Post("/house/create", api.NewHouseCreate(houseManager).Handle)

	if err := app.Listen(c.ServerAddress); err != nil {
		log.Fatal(err)
	}
}

func contextMiddleware(ctx context.Context) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		c.SetUserContext(ctx)
		return c.Next()
	}
}

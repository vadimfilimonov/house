package main

import (
	"context"
	"fmt"
	"log"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"

	"github.com/vadimfilimonov/house/internal/api"
	"github.com/vadimfilimonov/house/internal/api/flatcreate"
	"github.com/vadimfilimonov/house/internal/api/flatupdate"
	"github.com/vadimfilimonov/house/internal/api/housecreate"
	"github.com/vadimfilimonov/house/internal/api/houseget"
	"github.com/vadimfilimonov/house/internal/api/housesubscribe"
	"github.com/vadimfilimonov/house/internal/api/login"
	"github.com/vadimfilimonov/house/internal/api/register"
	"github.com/vadimfilimonov/house/internal/service/auth_token"
	"github.com/vadimfilimonov/house/internal/service/config"
	"github.com/vadimfilimonov/house/internal/service/flat"
	"github.com/vadimfilimonov/house/internal/service/house"
	"github.com/vadimfilimonov/house/internal/service/subscription"
	"github.com/vadimfilimonov/house/internal/service/user"
	"github.com/vadimfilimonov/house/internal/storage/pg"
	"github.com/vadimfilimonov/house/internal/storage/redis"
	flatStore "github.com/vadimfilimonov/house/internal/store/flat"
	houseStore "github.com/vadimfilimonov/house/internal/store/house"
	subscriptionStore "github.com/vadimfilimonov/house/internal/store/subscription"
	tokenStore "github.com/vadimfilimonov/house/internal/store/token"
	userStore "github.com/vadimfilimonov/house/internal/store/user"
)

func main() {
	ctx := context.Background()

	c := config.New()
	if err := c.Parse(); err != nil {
		log.Fatal(err)
	}

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.DatabaseHost, c.DatabasePort, c.PostgresUser, c.PostgresPassword, c.PostgresDatabaseName,
	)

	database, err := pg.New(connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	redisClient, err := redis.New(ctx, c.RedisHost, c.RedisPort, c.RedisPassword)
	if err != nil {
		log.Fatal(err)
	}
	defer redisClient.Close()

	uStore := userStore.New(database)
	hStore := houseStore.New(database)
	fStore := flatStore.New(database)
	sStore := subscriptionStore.New(database)
	tStore := tokenStore.New(redisClient)

	tokenManager := auth_token.NewToken([]byte(c.JwtSecretKey))
	userManager := user.New(uStore, tStore, tokenManager)
	houseManager := house.New(hStore)
	flatManager := flat.New(fStore)
	subscriptionManager := subscription.New(sStore)

	app := fiber.New()
	app.Use(contextMiddleware(ctx))

	publicGroup := app.Group("")
	publicGroup.Post("/login", login.New(userManager).Handle)
	publicGroup.Post("/register", register.New(userManager).Handle)

	authorizedGroup := app.Group("")
	authorizedGroup.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{
			Key: []byte(c.JwtSecretKey),
		},
		ContextKey: api.ContextKeyUser,
	}))
	authorizedGroup.Post("/house/create", housecreate.New(houseManager).Handle)
	authorizedGroup.Get("/house/:id", houseget.New(flatManager).Handle)
	authorizedGroup.Post("/house/:id/subscribe", housesubscribe.New(subscriptionManager).Handle)
	authorizedGroup.Post("/flat/create", flatcreate.New(flatManager).Handle)
	authorizedGroup.Post("/flat/update", flatupdate.New(flatManager).Handle)

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

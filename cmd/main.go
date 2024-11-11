package main

import (
	"log"

	"github.com/gofiber/fiber/v2"

	"github.com/vadimfilimonov/house/internal/api"
	"github.com/vadimfilimonov/house/internal/service/auth_token"
	"github.com/vadimfilimonov/house/internal/service/config"
	"github.com/vadimfilimonov/house/internal/service/user"
	tokenStorage "github.com/vadimfilimonov/house/internal/storage/token"
	userStorage "github.com/vadimfilimonov/house/internal/storage/user"
)

func main() {
	c := config.New()
	if err := c.Parse(); err != nil {
		log.Fatal(err)
	}

	uStorage, err := userStorage.New(c.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}

	defer uStorage.Close()

	tStorage, err := tokenStorage.New(c.RedisAddress, c.RedisPassword)
	if err != nil {
		log.Fatal(err)
	}

	defer tStorage.Close()

	tokenManager := auth_token.NewToken(c.JwtSecretKey)
	userManager := user.New(uStorage, tStorage, tokenManager)

	webApp := fiber.New()
	webApp.Post("/dummyLogin", api.NewDummyLogin(tokenManager, tStorage).Handle)
	webApp.Post("/login", api.NewLogin(userManager).Handle)
	webApp.Post("/register", api.NewRegister(userManager).Handle)

	if err := webApp.Listen(c.ServerAddress); err != nil {
		log.Fatal(err)
	}
}

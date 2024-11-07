package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/vadimfilimonov/house/internal/api"
	"github.com/vadimfilimonov/house/internal/service/auth_token"
	"github.com/vadimfilimonov/house/internal/service/config"
	"github.com/vadimfilimonov/house/internal/service/user"
	tokenStorage "github.com/vadimfilimonov/house/internal/storage/token"
	userStorage "github.com/vadimfilimonov/house/internal/storage/user"
)

func main() {
	c := config.New()
	err := c.Parse()
	if err != nil {
		log.Fatal(err)
	}

	uStorage, err := userStorage.New(c.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}

	tStorage, err := tokenStorage.New(c.RedisAddress, c.RedisPassword)
	if err != nil {
		log.Fatal(err)
	}

	tokenManager := auth_token.NewToken(c.JwtSecretKey, tStorage)
	userManager := user.New(uStorage, tStorage, tokenManager)

	r := chi.NewRouter()
	r.Post("/dummyLogin", api.NewDummyLogin(tokenManager, tStorage))
	r.Post("/login", api.NewLogin(userManager))
	r.Post("/register", api.NewRegister(userManager))

	err = http.ListenAndServe(c.ServerAddress, r)
	if err != nil {
		log.Fatal(err)
	}
}

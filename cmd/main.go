package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/vadimfilimonov/house/internal/api"
	"github.com/vadimfilimonov/house/internal/service/auth_token"
	"github.com/vadimfilimonov/house/internal/service/config"
	"github.com/vadimfilimonov/house/internal/service/user"
	"github.com/vadimfilimonov/house/internal/storage"
	userStorage "github.com/vadimfilimonov/house/internal/storage/user"
)

func main() {
	c := config.New()
	err := c.Parse()
	if err != nil {
		log.Fatal(err)
	}

	st, err := userStorage.GetStorage(storage.StorageTypeDatabase, c.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}

	token := auth_token.NewToken(c.JwtSecretKey)
	userManager := user.New(st)

	r := chi.NewRouter()
	r.Post("/dummyLogin", api.NewDummyLogin(token))
	r.Post("/register", api.NewRegister(userManager))

	err = http.ListenAndServe(c.ServerAddress, r)
	if err != nil {
		log.Fatal(err)
	}
}

package api

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	manager "github.com/vadimfilimonov/house/internal/service/user"
	storage "github.com/vadimfilimonov/house/internal/storage/user"
)

type LoginInput struct {
	ID       string `json:"user_id"`
	Password string `json:"password"`
}

type LoginOutput struct {
	Token string `json:"token"`
}

func NewLogin(userManager userManager) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()

		body, err := io.ReadAll(r.Body)
		defer r.Body.Close()

		if err != nil {
			BadRequest(&w, err.Error())
			return
		}

		var requestBody LoginInput

		err = json.Unmarshal([]byte(body), &requestBody)
		if err != nil {
			BadRequest(&w, err.Error())
			return
		}

		token, err := userManager.Login(ctx, requestBody.ID, requestBody.Password)
		if err != nil {
			if errors.Is(err, storage.ErrUserNotFound) {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}

			if errors.Is(err, manager.ErrWrongPassword) {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			BadRequest(&w, err.Error())
			return
		}

		response, err := json.Marshal(LoginOutput{Token: *token})
		if err != nil {
			BadRequest(&w, err.Error())
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(response)
	}
}

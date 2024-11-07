package api

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	manager "github.com/vadimfilimonov/house/internal/service/user"
)

type RegisterInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	UserType string `json:"user_type"`
}

type RegisterOutput struct {
	UserID string `json:"user_id"`
}

func NewRegister(userManager userManager) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		defer r.Body.Close()

		if err != nil {
			BadRequest(&w, err.Error())
			return
		}

		var requestBody RegisterInput

		err = json.Unmarshal([]byte(body), &requestBody)
		if err != nil {
			BadRequest(&w, err.Error())
			return
		}

		userID, err := userManager.Register(requestBody.Email, requestBody.Password, requestBody.UserType)
		if err != nil {
			if errors.Is(err, manager.ErrIncorrectInput) {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			BadRequest(&w, err.Error())
			return
		}

		if userID == nil {
			http.Error(w, "userID is empty", http.StatusBadRequest)
			return
		}

		response, err := json.Marshal(RegisterOutput{UserID: *userID})
		if err != nil {
			BadRequest(&w, err.Error())
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(response)
	}
}

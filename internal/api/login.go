package api

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
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
		body, err := io.ReadAll(r.Body)
		defer r.Body.Close()

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var requestBody LoginInput

		err = json.Unmarshal([]byte(body), &requestBody)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		token, err := userManager.Login(context.Background(), requestBody.ID, requestBody.Password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		response, err := json.Marshal(LoginOutput{Token: *token})
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(response)
	}
}

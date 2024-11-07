package api

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/vadimfilimonov/house/internal/service/auth_token"
)

const fakeEmail = "fake@mail.com"

type DummyLoginInput struct {
	UserType string `json:"user_type"`
}

type DummyLoginOutput struct {
	Token string `json:"token"`
}

func NewDummyLogin(tokenManager *auth_token.Token) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		defer r.Body.Close()

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var requestBody DummyLoginInput
		err = json.Unmarshal([]byte(body), &requestBody)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		token, err := tokenManager.Encode(fakeEmail)
		if err != nil || token == nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = tokenManager.Save(*token)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		response, err := json.Marshal(DummyLoginOutput{Token: *token})
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(response)
	}
}

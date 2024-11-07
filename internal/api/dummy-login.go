package api

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

const fakeEmail = "dummyLogin@mail.com"

type tokenManager interface {
	Encode(sub, userType string) (*string, error)
}

type tokenStorage interface {
	Add(ctx context.Context, key string, expiration time.Duration) error
}

type DummyLoginInput struct {
	UserType string `json:"user_type"`
}

type DummyLoginOutput struct {
	Token string `json:"token"`
}

func NewDummyLogin(tokenManager tokenManager, tokenStorage tokenStorage) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()

		body, err := io.ReadAll(r.Body)
		defer r.Body.Close()

		if err != nil {
			BadRequest(&w, err.Error())
			return
		}

		var requestBody DummyLoginInput
		err = json.Unmarshal([]byte(body), &requestBody)
		if err != nil {
			BadRequest(&w, err.Error())
			return
		}

		token, err := tokenManager.Encode(fakeEmail, requestBody.UserType)
		if err != nil || token == nil {
			BadRequest(&w, err.Error())
			return
		}

		err = tokenStorage.Add(ctx, *token, 24*time.Hour)
		if err != nil {
			BadRequest(&w, err.Error())
			return
		}

		response, err := json.Marshal(DummyLoginOutput{Token: *token})
		if err != nil {
			BadRequest(&w, err.Error())
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(response)
	}
}

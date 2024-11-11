package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/vadimfilimonov/house/internal/models"
)

type tokenManager interface {
	Encode(sub, userType string) (*string, error)
}

type tokenStorage interface {
	Add(ctx context.Context, key, value string, expiration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
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
		if err := json.Unmarshal([]byte(body), &requestBody); err != nil {
			BadRequest(&w, err.Error())
			return
		}

		var userID string
		switch requestBody.UserType {
		case models.UserTypeClient:
			userID = models.FakeClientUserID
		case models.UserTypeModerator:
			userID = models.FakeModeratorUserID
		default:
			http.Error(w, fmt.Sprintf("user type %s is not supported", requestBody.UserType), http.StatusNotFound)
			return
		}

		var token *string

		savedToken, err := tokenStorage.Get(ctx, userID)
		if err == nil {
			token = &savedToken
		} else {
			token, err = tokenManager.Encode(userID, requestBody.UserType)
			if err != nil {
				BadRequest(&w, err.Error())
				return
			} else if token == nil {
				BadRequest(&w, "token is nil")
				return
			}

			if err := tokenStorage.Add(ctx, userID, *token, 24*time.Hour); err != nil {
				BadRequest(&w, err.Error())
				return
			}
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

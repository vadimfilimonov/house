package auth_token

import (
	"context"
	"fmt"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

type storage interface {
	Add(ctx context.Context, key, value string, expiration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
}

type Token struct {
	secretKey string
	storage   storage
}

func NewToken(secretKey string, storage storage) *Token {
	return &Token{
		secretKey: secretKey,
		storage:   storage,
	}
}

func (t *Token) Encode(sub string) (*string, error) {
	payload := jwt.MapClaims{
		"sub": sub,
		"exp": time.Now().Add(time.Hour * 72).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	signedToken, err := token.SignedString([]byte(t.secretKey))
	if err != nil {
		return nil, fmt.Errorf("cannot sign token: %w", err)
	}

	return &signedToken, nil
}

func (t *Token) Decode(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return t.secretKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("cannot parse token: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("cannot parse claims")
	}

	return claims, nil
}

// Save saves token to storage
func (t *Token) Save(ctx context.Context, key, tokenString string) error {
	err := t.storage.Add(ctx, key, tokenString, 24*time.Hour)
	return err
}

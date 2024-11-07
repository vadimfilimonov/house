package auth_token

import (
	"context"
	"fmt"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

const (
	claimsKeySub      = "sub"
	claimsKeyUserType = "userType"
)

type storage interface {
	Add(ctx context.Context, key string, expiration time.Duration) error
	Has(ctx context.Context, key string) bool
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

func (t *Token) Encode(sub, userType string) (*string, error) {
	payload := jwt.MapClaims{
		claimsKeySub:      sub,
		claimsKeyUserType: userType,
		"exp":             time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	signedToken, err := token.SignedString([]byte(t.secretKey))
	if err != nil {
		return nil, fmt.Errorf("cannot sign token: %w", err)
	}

	return &signedToken, nil
}

func (t *Token) Decode(tokenString string) (id *string, userType *string, err error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return t.secretKey, nil
	})

	if err != nil {
		return nil, nil, fmt.Errorf("cannot parse token: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, nil, fmt.Errorf("cannot parse claims")
	}

	sub, exists := claims[claimsKeySub]
	if exists {
		userID, ok := sub.(string)
		if ok {
			id = &userID
		}
	}

	userTypeValue, exists := claims[claimsKeyUserType]
	if exists {
		userTypeStr, ok := userTypeValue.(string)
		if ok {
			userType = &userTypeStr
		}
	}

	return id, userType, nil
}

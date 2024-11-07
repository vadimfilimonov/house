package auth_token

import (
	"fmt"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

type Token struct {
	secretKey string
}

func NewToken(secretKey string) *Token {
	return &Token{
		secretKey: secretKey,
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
func (t *Token) Save(tokenString string) error {
	return nil
}

package api

import "context"

type userManager interface {
	Register(email, password, userType string) (userID *string, err error)
	Login(ctx context.Context, id string, password string) (token *string, err error)
}

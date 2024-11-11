package api

import "context"

type userManager interface {
	Register(ctx context.Context, email, password, userType string) (userID *string, err error)
	Login(ctx context.Context, id string, password string) (token *string, err error)
}

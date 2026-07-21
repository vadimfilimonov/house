package login

import "context"

type userManager interface {
	// Login authenticates a user and returns an authorization token.
	Login(ctx context.Context, id string, password string) (token *string, err error)
}

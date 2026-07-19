package register

import "context"

type userManager interface {
	// Register creates a user with the requested role.
	Register(ctx context.Context, email, password, userType string) (userID *string, err error)
}

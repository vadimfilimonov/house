package login

import "github.com/vadimfilimonov/house/internal/api"

type Input struct {
	// User email.
	Email string `json:"email"`
	// User password.
	Password string `json:"password"`
}

// Validate checks that the login request matches API constraints.
func (i Input) Validate() error {
	if err := api.ValidateEmail(i.Email); err != nil {
		return err
	}

	if err := api.ValidatePassword(i.Password); err != nil {
		return err
	}

	return nil
}

type Output struct {
	// Authorization token.
	Token string `json:"token"`
}

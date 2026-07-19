package register

import "github.com/vadimfilimonov/house/internal/api"

type Input struct {
	// User email.
	Email string `json:"email"`
	// User password.
	Password string `json:"password"`
	// User type: client or moderator.
	UserType string `json:"user_type"`
}

// Validate checks that the registration request matches API constraints.
func (i Input) Validate() error {
	if err := api.ValidateEmail(i.Email); err != nil {
		return err
	}

	if err := api.ValidatePassword(i.Password); err != nil {
		return err
	}

	if err := api.ValidateUserType(i.UserType); err != nil {
		return err
	}

	return nil
}

type Output struct {
	// Created user identifier.
	UserID string `json:"user_id"`
}

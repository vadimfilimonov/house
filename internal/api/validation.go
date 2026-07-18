package api

import (
	"fmt"
	"regexp"
	"unicode/utf8"
)

const (
	minPasswordSize = 8

	userTypeClient    = "client"
	userTypeModerator = "moderator"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func validateEmail(email string) error {
	if !emailRegex.MatchString(email) {
		return fmt.Errorf("email is not valid")
	}

	return nil
}

func validatePassword(password string) error {
	if utf8.RuneCountInString(password) < minPasswordSize {
		return fmt.Errorf("password is too short, use min %d symbols", minPasswordSize)
	}

	return nil
}

func validateUserType(userType string) error {
	if userType != userTypeClient && userType != userTypeModerator {
		return fmt.Errorf("user_type is incorrect")
	}

	return nil
}

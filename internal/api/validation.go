package api

import (
	"fmt"
	"regexp"
	"unicode/utf8"
)

const (
	minPasswordSize = 8
	maxStringSize   = 255

	userTypeClient    = "client"
	userTypeModerator = "moderator"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// ValidateEmail checks that an email matches API constraints.
func ValidateEmail(email string) error {
	if err := ValidateMaxStringSize("email", email); err != nil {
		return err
	}

	if !emailRegex.MatchString(email) {
		return fmt.Errorf("email is not valid")
	}

	return nil
}

// ValidatePassword checks that a password matches API constraints.
func ValidatePassword(password string) error {
	if utf8.RuneCountInString(password) < minPasswordSize {
		return fmt.Errorf("password is too short, use min %d symbols", minPasswordSize)
	}

	return nil
}

// ValidateUserType checks that a user type is supported by the API.
func ValidateUserType(userType string) error {
	if userType != userTypeClient && userType != userTypeModerator {
		return fmt.Errorf("user_type is incorrect")
	}

	return nil
}

// ValidateMaxStringSize checks that a string field fits into database limits.
func ValidateMaxStringSize(field, value string) error {
	if utf8.RuneCountInString(value) > maxStringSize {
		return fmt.Errorf("%s cannot be longer than %d symbols", field, maxStringSize)
	}

	return nil
}

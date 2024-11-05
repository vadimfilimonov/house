package user

import (
	"fmt"
	"regexp"
	"unicode/utf8"
)

const (
	userTypeClient    = "client"
	userTypeModerator = "moderator"
)

const (
	minPasswordSize = 8
)

func (u *UserManager) validate(email, password, userType string) error {
	if err := u.validateUserType(userType); err != nil {
		return err
	}

	if err := u.validatePassword(password); err != nil {
		return err
	}

	if err := u.validateEmail((email)); err != nil {
		return err
	}

	return nil
}

func (u *UserManager) validateUserType(userType string) error {
	if userType != userTypeClient && userType != userTypeModerator {
		return fmt.Errorf("incorrect user type: %s", userType)
	}

	return nil
}

func (u *UserManager) validatePassword(password string) error {
	runeCount := utf8.RuneCountInString(password)
	if runeCount < minPasswordSize {
		return fmt.Errorf("password is too short, use min %d symbols", runeCount)
	}

	return nil
}

func (u *UserManager) validateEmail(email string) error {
	if !u.isValidEmail(email) {
		return fmt.Errorf("email %s is not valid", email)
	}

	return nil
}

// isValidEmail checks if the string is a valid email address
func (u *UserManager) isValidEmail(email string) bool {
	const emailRegex = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`

	re := regexp.MustCompile(emailRegex)

	return re.MatchString(email)
}

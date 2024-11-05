package storage

import (
	"fmt"

	"github.com/vadimfilimonov/house/internal/models"
)

type Memory struct {
	users map[string]models.User
}

func NewMemory() *Memory {
	users := make(map[string]models.User)

	return &Memory{users: users}
}

// Add adds user to memory
func (m *Memory) Add(email, hashedPassword, userType string) (*string, error) {
	userID, err := generateUserID()
	if err != nil {
		return nil, fmt.Errorf("cannot add user to memory storage: %w", err)
	}

	if userID == nil {
		return nil, fmt.Errorf("userID was not generated")
	}

	_, isEmailBusy := m.getByEmail(email)
	if isEmailBusy {
		return nil, fmt.Errorf("email \"%s\" is already used", email)
	}

	m.users[*userID] = models.User{
		ID:       *userID,
		Email:    email,
		Password: hashedPassword,
		UserType: userType,
	}

	return userID, nil
}

// Get gets user from memory
func (m *Memory) Get(id string) (*models.User, error) {
	user, ok := m.users[id]
	if !ok {
		return nil, fmt.Errorf("user with id %s is not found", id)
	}

	return &user, nil
}

func (s *Memory) getByEmail(email string) (*models.User, bool) {
	for _, value := range s.users {
		if email == value.Email {
			return &value, true
		}
	}

	return nil, false
}

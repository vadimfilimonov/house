package models

const (
	UserTypeClient    = "client"
	UserTypeModerator = "moderator"
)

const (
	FakeClientUserID    = "fake-client"
	FakeModeratorUserID = "fake-moderator"
)

type User struct {
	ID       string `json:"user_id"   db:"user_id"   example:"123e4567-e89b-12d3-a456-426655440000"`
	Email    string `json:"email"     db:"email"     example:"ivanov@mail.ru"`
	Password string `json:"password"  db:"password"  example:"qwerty"`
	UserType string `json:"user_type" db:"user_type" example:"client"`
}

func (u *User) IsClient() bool {
	return u.UserType == UserTypeClient
}

func (u *User) IsModerator() bool {
	return u.UserType == UserTypeModerator
}

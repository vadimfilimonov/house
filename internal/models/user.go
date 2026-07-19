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
	ID       string `db:"user_id"   example:"123e4567-e89b-12d3-a456-426655440000"` // User identifier.
	Email    string `db:"email"     example:"ivanov@mail.ru"`                       // User email.
	Password string `db:"password"  example:"qwerty"`                               // User password hash.
	UserType string `db:"user_type" example:"client"`                               // User type.
}

func (u *User) IsClient() bool {
	return u.UserType == UserTypeClient
}

func (u *User) IsModerator() bool {
	return u.UserType == UserTypeModerator
}

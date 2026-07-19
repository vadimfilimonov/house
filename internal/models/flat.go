package models

type Status string // Flat moderation status.

func (s Status) String() string {
	return string(s)
}

const (
	CreatedStatus      Status = "created"
	ApprovedStatus     Status = "approved"
	DeclinedStatus     Status = "declined"
	OnModerationStatus Status = "on moderation"
)

type Flat struct {
	ID      int     `db:"id"       example:"123456"`  // Flat identifier.
	Number  int     `db:"number"   example:"12345"`   // Flat number inside the house.
	HouseID HouseID `db:"house_id" example:"12345"`   // House identifier linked to the flat.
	Price   int     `db:"price"    example:"10000"`   // Flat price in conventional units.
	Rooms   int     `db:"rooms"    example:"4"`       // Number of rooms in the flat.
	Status  Status  `db:"status"   example:"created"` // Flat moderation status.
}

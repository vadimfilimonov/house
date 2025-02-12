package models

type Status string // Статус квартиры

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
	Number  int     `db:"number"   example:"12345"` // Номер квартиры
	HouseID HouseID `db:"house_id" example:"12345"`
	Price   int     `db:"price"    example:"10000"` // Цена квартиры в у.е.
	Rooms   int     `db:"rooms"    example:"4"`     // Количество комнат в квартире
	Status  Status  `db:"status"   example:"created"`
}

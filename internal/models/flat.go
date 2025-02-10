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
	ID      int     `json:"id"       db:"id"       example:"12345"` // Номер квартиры
	HouseID HouseID `json:"house_id" db:"house_id" example:"12345"`
	Price   int     `json:"price"    db:"price"    example:"10000"` // Цена квартиры в у.е.
	Rooms   int     `json:"rooms"    db:"rooms"    example:"4"`     // Количество комнат в квартире
	Status  Status  `json:"status"   db:"status"   example:"created"`
}

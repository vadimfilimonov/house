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
	ID      int     `db:"id"       example:"123456"`  // Идентификатор квартиры
	Number  int     `db:"number"   example:"12345"`   // Номер квартиры в доме
	HouseID HouseID `db:"house_id" example:"12345"`   // Идентификатор дома, к которому относится квартира
	Price   int     `db:"price"    example:"10000"`   // Цена квартиры в у.е.
	Rooms   int     `db:"rooms"    example:"4"`       // Количество комнат в квартире
	Status  Status  `db:"status"   example:"created"` // Статус модерации квартиры
}

package models

type HouseID int // House identifier.

func (h HouseID) Int() int {
	return int(h)
}

type House struct {
	ID        HouseID `db:"id"         example:"12345"`                           // House identifier.
	Address   string  `db:"address"    example:"Лесная улица, 7, Москва, 125196"` // House address.
	Year      int     `db:"year"       example:"2000"`                            // House construction year.
	Developer *string `db:"developer"  example:"Мэрия города"`                    // House developer.
	CreatedAt *string `db:"created_at" example:"2017-07-21T17:32:28Z"`            // House creation date.
	UpdateAt  *string `db:"update_at"  example:"2017-07-21T17:32:28Z"`            // Date when a flat was last added to the house.
}

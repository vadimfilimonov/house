package models

type HouseID int // Уникальный номер дома

func (h HouseID) Int() int {
	return int(h)
}

type House struct {
	ID        HouseID `db:"id"         example:"12345"`                           // Идентификатор дома
	Address   string  `db:"address"    example:"Лесная улица, 7, Москва, 125196"` // Адрес дома
	Year      int     `db:"year"       example:"2000"`                            // Год постройки дома
	Developer *string `db:"developer"  example:"Мэрия города"`                    // Застройщик дома
	CreatedAt *string `db:"created_at" example:"2017-07-21T17:32:28Z"`            // Дата создания дома в базе
	UpdateAt  *string `db:"update_at"  example:"2017-07-21T17:32:28Z"`            // Дата последнего добавления новой квартиры в дом
}

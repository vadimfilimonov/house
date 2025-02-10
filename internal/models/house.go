package models

type HouseID int // Уникальный номер дома

func (h HouseID) Int() int {
	return int(h)
}

type House struct {
	ID        HouseID `db:"id"         example:"12345"`
	Address   string  `db:"address"    example:"Лесная улица, 7, Москва, 125196"` // Адрес
	Year      int     `db:"year"       example:"2000"`                            // Год постройки
	Developer *string `db:"developer"  example:"Мэрия города"`                    // Застройщик
	CreatedAt *string `db:"created_at" example:"2017-07-21T17:32:28Z"`            // Дата создания дома в базе
	UpdateAt  *string `db:"update_at"  example:"2017-07-21T17:32:28Z"`            // Дата последнего добавления новой квартиры дома
}

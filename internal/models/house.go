package models

type House struct {
	ID        int     `json:"id"                   db:"id"         example:"12345"`                           // Уникальный номер дома
	Address   string  `json:"address"              db:"address"    example:"Лесная улица, 7, Москва, 125196"` // Адрес
	Year      int     `json:"year"                 db:"year"       example:"2000"`                            // Год постройки
	Developer *string `json:"developer,omitempty"  db:"developer"  example:"Мэрия города"`                    // Застройщик
	CreatedAt *string `json:"created_at,omitempty" db:"created_at" example:"2017-07-21T17:32:28Z"`            // Дата создания дома в базе
	UpdateAt  *string `json:"update_at,omitempty"  db:"update_at"  example:"2017-07-21T17:32:28Z"`            // Дата последнего добавления новой квартиры дома
}

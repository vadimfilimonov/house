package models

type House struct {
	ID        string  `json:"id"                   db:"id"`         // Уникальный номер дома
	Address   string  `json:"address"              db:"address"`    // Адрес
	Year      int     `json:"year"                 db:"year"`       // Год постройки
	Developer *string `json:"developer,omitempty"  db:"developer"`  // Застройщик
	CreatedAt *string `json:"created_at,omitempty" db:"created_at"` // Дата создания дома в базе
	UpdateAt  *string `json:"update_at,omitempty"  db:"update_at"`  // Дата последнего добавления новой квартиры дома
}

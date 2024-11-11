package models

type House struct {
	ID        string  `json:"id"                   db:"id"`
	Address   string  `json:"address"              db:"address"`
	Year      int     `json:"year"                 db:"year"`
	Developer *string `json:"developer,omitempty"  db:"developer"`
	CreatedAt *string `json:"created_at,omitempty" db:"created_at"`
	UpdateAt  *string `json:"update_at,omitempty"  db:"update_at"`
}

package api

type FlatOutput struct {
	// Flat identifier.
	ID int `json:"id"`
	// Flat number inside the house.
	Number int `json:"number"`
	// House identifier linked to the flat.
	HouseID int `json:"house_id"`
	// Flat price in conventional units.
	Price int `json:"price"`
	// Number of rooms in the flat.
	Rooms int `json:"rooms"`
	// Flat moderation status.
	Status string `json:"status"`
}

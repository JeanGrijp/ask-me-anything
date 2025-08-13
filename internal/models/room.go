package models

type Room struct {
	ID      int64  `json:"id" db:"id"`
	Name    string `json:"name" db:"name"`
	OwnerID int64  `json:"owner_id" db:"owner_id"`
}

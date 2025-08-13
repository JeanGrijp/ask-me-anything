package models

type Answer struct {
	ID       int64    `json:"id" db:"id"`
	Question Question `json:"question" db:"question"`
	Answer   string   `json:"answer" db:"answer"`
	UserID   int64    `json:"user_id" db:"user_id"`
}

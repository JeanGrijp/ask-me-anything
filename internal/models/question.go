package models

type Question struct {
	ID         int64  `json:"id" db:"id"`
	Content    string `json:"content" db:"content"`
	UserID     int64  `json:"user_id" db:"user_id"`
	LikeCount  int64  `json:"like_count" db:"like_count"`
	IsAnswered bool   `json:"is_answered" db:"is_answered"`
}

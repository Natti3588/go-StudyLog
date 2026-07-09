package domain

import "time"

type StudyLog struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	CategoryID  string    `json:"category_id"`
	StudiedOn   time.Time `json:"studied_on"`
	DurationMin int       `json:"duration_min"`
	Memo        string    `json:"memo,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

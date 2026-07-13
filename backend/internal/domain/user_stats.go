package domain

import "time"

type UserStats struct {
	UserID        string     `json:"user_id"`
	TotalMin      int        `json:"total_min"`
	CurrentStreak int        `json:"current_streak"`
	LongestStreak int        `json:"longest_streak"`
	LastStudiedOn *time.Time `json:"last_studied_on"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

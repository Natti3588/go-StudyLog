package domain

import "time"

type WeeklyGoal struct {
	UserID    string    `json:"user_id"`
	WeekStart time.Time `json:"week_start"`
	TargetMin int       `json:"target_min"`
}

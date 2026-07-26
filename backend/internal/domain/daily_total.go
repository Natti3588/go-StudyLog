package domain

import "time"

type DailyTotal struct {
	Date     time.Time `json:"date"`
	TotalMin int       `json:"total_min"`
}

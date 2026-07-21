package domain

import "time"

type DailyTotal struct {
	Date     time.Time
	TotalMin int
}

package domain

type StatsSummary struct {
	TotalMin        int `json:"total_min"`
	CurrentStreak   int `json:"current_streak"`
	LongestStreak   int `json:"longest_streak"`
	WeeklyTargetMin int `json:"weekly_target_min"`
	WeeklyActualMin int `json:"weekly_actual_min"`
}

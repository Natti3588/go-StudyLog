export interface Category {
  id: string
  user_id: string | null
  name: string
  created_at: string
}

export interface StudyLog {
  id: string
  user_id: string
  category_id: string
  studied_on: string
  duration_min: number
  memo?: string
  created_at: string
  updated_at: string
}

export interface User {
  id: string
  email: string
  is_admin: boolean
  created_at: string
}

export interface WeeklyGoal {
  user_id: string
  week_start: string
  target_min: number
}

export interface UserStats {
  user_id: string
  total_min: number
  current_streak: number
  longest_streak: number
  last_studied_on: string | null
  updated_at: string
}

export interface DailyTotal {
  date: string
  total_min: number
}

export interface StatsSummary {
  total_min: number
  current_streak: number
  longest_streak: number
  weekly_target_min: number
  weekly_actual_min: number
}

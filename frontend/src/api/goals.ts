import { apiRequest } from './client'
import type { WeeklyGoal } from '../types/domain'

export interface WeeklyGoalInput {
  week_start: string
  target_min: number
}

export function setWeeklyGoal(input: WeeklyGoalInput): Promise<WeeklyGoal> {
  return apiRequest<WeeklyGoal>('/goals/weekly', { method: 'PUT', body: input })
}

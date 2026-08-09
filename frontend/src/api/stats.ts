import { apiRequest } from './client'
import type { StatsSummary, DailyTotal } from '../types/domain'

export function getStatsSummary(): Promise<StatsSummary> {
  return apiRequest<StatsSummary>('/stats/summary')
}

export function getHeatmap(year: number): Promise<DailyTotal[]> {
  return apiRequest<DailyTotal[]>(`/stats/heatmap?year=${year}`)
}

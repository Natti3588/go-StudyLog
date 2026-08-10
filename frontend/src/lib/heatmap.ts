import type { DailyTotal } from '../types/domain'

export interface HeatmapDatum {
  date: string
  count: number
  level: number
}

export function levelForMinutes(totalMin: number): number {
  if (totalMin <= 0) return 0
  if (totalMin < 30) return 1
  if (totalMin < 60) return 2
  if (totalMin < 120) return 3
  return 4
}

export function buildHeatmapData(year: number, totals: DailyTotal[]): HeatmapDatum[] {
  const byDate = new Map<string, number>()
  for (const t of totals) {
    byDate.set(t.date.slice(0, 10), t.total_min)
  }

  const data: HeatmapDatum[] = []
  const start = new Date(Date.UTC(year, 0, 1))
  const end = new Date(Date.UTC(year, 11, 31))

  for (let d = start; d.getTime() <= end.getTime(); d = new Date(d.getTime() + 86_400_000)) {
    const dateStr = d.toISOString().slice(0, 10)
    const totalMin = byDate.get(dateStr) ?? 0
    data.push({ date: dateStr, count: totalMin, level: levelForMinutes(totalMin) })
  }

  return data
}

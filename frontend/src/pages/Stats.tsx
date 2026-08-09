import { useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import { ActivityCalendar } from 'react-activity-calendar'
import { getStatsSummary, getHeatmap } from '../api/stats'
import { buildHeatmapData } from '../lib/heatmap'
import { LoadingState } from '../components/LoadingState'
import { ErrorState } from '../components/ErrorState'

function availableYears(): number[] {
  const current = new Date().getUTCFullYear()
  return [current, current - 1, current - 2]
}

export function Stats() {
  const [year, setYear] = useState(() => new Date().getUTCFullYear())

  const summaryQuery = useQuery({ queryKey: ['stats-summary'], queryFn: getStatsSummary })
  const heatmapQuery = useQuery({
    queryKey: ['stats-heatmap', year],
    queryFn: () => getHeatmap(year),
  })

  if (summaryQuery.isPending || heatmapQuery.isPending) {
    return <LoadingState />
  }

  if (summaryQuery.isError || heatmapQuery.isError) {
    return (
      <ErrorState
        onRetry={() => {
          summaryQuery.refetch()
          heatmapQuery.refetch()
        }}
      />
    )
  }

  const summary = summaryQuery.data
  const heatmapData = buildHeatmapData(year, heatmapQuery.data)

  return (
    <div className="p-6">
      <h1 className="text-xl font-bold text-ink mb-4">統計</h1>

      <div className="grid grid-cols-1 sm:grid-cols-3 gap-4 mb-6">
        <div className="bg-surface rounded-lg p-4">
          <p className="text-sm text-ink-muted">累計学習時間</p>
          <p className="text-2xl font-bold text-ink">{summary.total_min}分</p>
        </div>
        <div className="bg-surface rounded-lg p-4">
          <p className="text-sm text-ink-muted">現在のストリーク</p>
          <p className="text-2xl font-bold text-ink">{summary.current_streak}日</p>
        </div>
        <div className="bg-surface rounded-lg p-4">
          <p className="text-sm text-ink-muted">最長ストリーク</p>
          <p className="text-2xl font-bold text-ink">{summary.longest_streak}日</p>
        </div>
      </div>

      <div className="bg-surface rounded-lg p-4">
        <label className="block mb-3">
          <span className="block text-sm text-ink-secondary mb-1">年</span>
          <select
            value={year}
            onChange={(e) => setYear(Number(e.target.value))}
            className="border border-hairline rounded-xs px-2 py-1.5"
          >
            {availableYears().map((y) => (
              <option key={y} value={y}>
                {y}年
              </option>
            ))}
          </select>
        </label>
        <ActivityCalendar data={heatmapData} />
      </div>
    </div>
  )
}

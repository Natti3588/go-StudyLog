import { useState } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { useForm } from 'react-hook-form'
import { Link } from 'react-router'
import { useAuth } from '../context/AuthContext'
import { getStatsSummary } from '../api/stats'
import { setWeeklyGoal } from '../api/goals'
import { listLogs } from '../api/logs'
import { LoadingState } from '../components/LoadingState'
import { ErrorState } from '../components/ErrorState'
import { localDateString } from '../lib/date'

function currentWeekStart(): string {
  const now = new Date()
  const day = now.getDay()
  const diffToMonday = day === 0 ? 6 : day - 1
  const monday = new Date(now.getFullYear(), now.getMonth(), now.getDate() - diffToMonday)
  return `${localDateString(monday)}T00:00:00Z`
}

interface GoalFormValues {
  target_min: number
}

export function Dashboard() {
  const { user, logout } = useAuth()
  const queryClient = useQueryClient()
  const [isEditingGoal, setIsEditingGoal] = useState(false)

  const summaryQuery = useQuery({ queryKey: ['stats-summary'], queryFn: getStatsSummary })
  const logsQuery = useQuery({ queryKey: ['logs'], queryFn: listLogs })

  const goalMutation = useMutation({
    mutationFn: (values: GoalFormValues) =>
      setWeeklyGoal({ week_start: currentWeekStart(), target_min: values.target_min }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['stats-summary'] })
      setIsEditingGoal(false)
    },
  })

  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
  } = useForm<GoalFormValues>()

  if (summaryQuery.isPending || logsQuery.isPending) {
    return <LoadingState />
  }

  if (summaryQuery.isError || logsQuery.isError) {
    return (
      <ErrorState
        onRetry={() => {
          summaryQuery.refetch()
          logsQuery.refetch()
        }}
      />
    )
  }

  const summary = summaryQuery.data
  const recentLogs = logsQuery.data.slice(0, 5)
  const progressPercent =
    summary.weekly_target_min > 0
      ? Math.min(100, Math.round((summary.weekly_actual_min / summary.weekly_target_min) * 100))
      : 0

  return (
    <div className="p-6">
      <h1 className="text-xl font-bold text-ink mb-1">ようこそ、{user?.email}</h1>
      <button onClick={() => logout()} className="text-sm text-ink-muted mb-6">
        ログアウト
      </button>

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

      <div className="bg-surface rounded-lg p-4 mb-6">
        <div className="flex items-center justify-between mb-2">
          <p className="text-sm text-ink-muted">
            週間目標: {summary.weekly_actual_min} / {summary.weekly_target_min}分
          </p>
          {!isEditingGoal && (
            <button onClick={() => setIsEditingGoal(true)} className="text-sm text-primary">
              目標を変更
            </button>
          )}
        </div>
        <div className="w-full bg-canvas-soft rounded-full h-2">
          <div className="bg-primary h-2 rounded-full" style={{ width: `${progressPercent}%` }} />
        </div>
        {isEditingGoal && (
          <form
            onSubmit={handleSubmit((values) => goalMutation.mutate(values))}
            className="mt-3 flex items-center gap-2"
          >
            <input
              type="number"
              className="border border-hairline rounded-xs px-2 py-1 w-24"
              {...register('target_min', {
                required: true,
                valueAsNumber: true,
                min: { value: 1, message: '1分以上を入力してください' },
              })}
            />
            <button
              type="submit"
              disabled={isSubmitting}
              className="bg-primary text-white rounded-md px-3 py-1 text-sm disabled:opacity-50"
            >
              保存
            </button>
            {errors.target_min && <p className="text-error text-sm">{errors.target_min.message}</p>}
          </form>
        )}
      </div>

      <div className="bg-surface rounded-lg p-4">
        <div className="flex items-center justify-between mb-3">
          <p className="text-sm text-ink-muted">最近の記録</p>
          <Link to="/logs" className="text-sm text-primary">
            すべて見る
          </Link>
        </div>
        {recentLogs.length === 0 ? (
          <p className="text-sm text-ink-faint">まだ記録がありません</p>
        ) : (
          <ul className="flex flex-col gap-2">
            {recentLogs.map((log) => (
              <li key={log.id} className="flex justify-between text-sm">
                <span className="text-ink">{log.studied_on.slice(0, 10)}</span>
                <span className="text-ink-muted">{log.duration_min}分</span>
              </li>
            ))}
          </ul>
        )}
      </div>
    </div>
  )
}

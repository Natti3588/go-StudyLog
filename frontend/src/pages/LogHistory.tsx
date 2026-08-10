import { useState } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { listCategories } from '../api/categories'
import { listLogs, updateLog, deleteLog } from '../api/logs'
import { LogForm, type LogFormValues } from '../components/LogForm'
import { Modal } from '../components/Modal'
import { LoadingState } from '../components/LoadingState'
import { ErrorState } from '../components/ErrorState'
import type { StudyLog } from '../types/domain'

export function LogHistory() {
  const queryClient = useQueryClient()
  const [editingLog, setEditingLog] = useState<StudyLog | null>(null)

  const logsQuery = useQuery({ queryKey: ['logs'], queryFn: listLogs })
  const categoriesQuery = useQuery({ queryKey: ['categories'], queryFn: listCategories })

  const updateMutation = useMutation({
    mutationFn: (values: LogFormValues) => updateLog(editingLog!.id, values),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['logs'] })
      queryClient.invalidateQueries({ queryKey: ['stats-summary'] })
      setEditingLog(null)
    },
  })

  const deleteMutation = useMutation({
    mutationFn: (id: string) => deleteLog(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['logs'] })
      queryClient.invalidateQueries({ queryKey: ['stats-summary'] })
    },
  })

  function handleDelete(id: string) {
    if (window.confirm('この記録を削除しますか?')) {
      deleteMutation.mutate(id)
    }
  }

  if (logsQuery.isPending || categoriesQuery.isPending) {
    return <LoadingState />
  }

  if (logsQuery.isError || categoriesQuery.isError) {
    return (
      <ErrorState
        onRetry={() => {
          logsQuery.refetch()
          categoriesQuery.refetch()
        }}
      />
    )
  }

  const categoryName = (categoryId: string) =>
    categoriesQuery.data.find((c) => c.id === categoryId)?.name ?? '不明'

  return (
    <div className="p-6">
      <h1 className="text-xl font-bold text-ink mb-4">ログ履歴</h1>
      <table className="w-full text-sm">
        <thead>
          <tr className="border-b border-hairline text-left text-ink-muted">
            <th className="py-2">日付</th>
            <th className="py-2">カテゴリ</th>
            <th className="py-2">時間</th>
            <th className="py-2">メモ</th>
            <th className="py-2"></th>
          </tr>
        </thead>
        <tbody>
          {logsQuery.data.map((log) => (
            <tr key={log.id} className="border-b border-hairline">
              <td className="py-2">{log.studied_on.slice(0, 10)}</td>
              <td className="py-2">{categoryName(log.category_id)}</td>
              <td className="py-2">{log.duration_min}分</td>
              <td className="py-2">{log.memo}</td>
              <td className="py-2 text-right">
                <button onClick={() => setEditingLog(log)} className="text-primary mr-3">
                  編集
                </button>
                <button onClick={() => handleDelete(log.id)} className="text-error">
                  削除
                </button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>

      {editingLog && (
        <Modal onClose={() => setEditingLog(null)}>
          <h2 className="text-lg font-bold text-ink mb-4">記録を編集</h2>
          <LogForm
            categories={categoriesQuery.data}
            defaultValues={{
              category_id: editingLog.category_id,
              studied_on: editingLog.studied_on.slice(0, 10),
              duration_min: editingLog.duration_min,
              memo: editingLog.memo ?? '',
            }}
            onSubmit={async (values) => {
              await updateMutation.mutateAsync(values)
            }}
            submitLabel="更新する"
          />
        </Modal>
      )}
    </div>
  )
}

import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { useNavigate } from 'react-router'
import { listCategories } from '../api/categories'
import { createLog } from '../api/logs'
import { LogForm, type LogFormValues } from '../components/LogForm'
import { LoadingState } from '../components/LoadingState'
import { ErrorState } from '../components/ErrorState'

export function RecordInput() {
  const navigate = useNavigate()
  const queryClient = useQueryClient()
  const categoriesQuery = useQuery({ queryKey: ['categories'], queryFn: listCategories })
  const createMutation = useMutation({
    mutationFn: createLog,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['logs'] })
      queryClient.invalidateQueries({ queryKey: ['stats-summary'] })
      navigate('/logs')
    },
  })

  if (categoriesQuery.isPending) {
    return <LoadingState />
  }

  if (categoriesQuery.isError) {
    return <ErrorState onRetry={() => categoriesQuery.refetch()} />
  }

  return (
    <div className="p-6 max-w-md">
      <h1 className="text-xl font-bold text-ink mb-4">記録入力</h1>
      <LogForm
        categories={categoriesQuery.data}
        onSubmit={async (values: LogFormValues) => {
          await createMutation.mutateAsync(values)
        }}
        submitLabel="記録する"
      />
    </div>
  )
}

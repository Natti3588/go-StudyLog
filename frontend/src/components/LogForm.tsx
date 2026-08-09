import { useForm } from 'react-hook-form'
import { useState } from 'react'
import type { Category } from '../types/domain'

export interface LogFormValues {
  category_id: string
  studied_on: string
  duration_min: number
  memo: string
}

interface LogFormProps {
  categories: Category[]
  defaultValues?: LogFormValues
  onSubmit: (values: LogFormValues) => Promise<void>
  submitLabel: string
}

function today(): string {
  return new Date().toISOString().slice(0, 10)
}

export function LogForm({ categories, defaultValues, onSubmit, submitLabel }: LogFormProps) {
  const [serverError, setServerError] = useState<string | null>(null)
  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
  } = useForm<LogFormValues>({
    defaultValues: defaultValues ?? {
      category_id: '',
      studied_on: today(),
      duration_min: 30,
      memo: '',
    },
  })

  async function submit(values: LogFormValues) {
    setServerError(null)
    try {
      await onSubmit(values)
    } catch {
      setServerError('エラーが発生しました。時間をおいて再度お試しください')
    }
  }

  return (
    <form onSubmit={handleSubmit(submit)} noValidate>
      <label className="block mb-3">
        <span className="block text-sm text-ink-secondary mb-1">カテゴリ</span>
        <select
          className="w-full border border-hairline rounded-xs px-2 py-1.5"
          {...register('category_id', { required: 'カテゴリを選択してください' })}
        >
          <option value="">選択してください</option>
          {categories.map((c) => (
            <option key={c.id} value={c.id}>
              {c.name}
            </option>
          ))}
        </select>
        {errors.category_id && <p className="text-error text-sm mt-1">{errors.category_id.message}</p>}
      </label>
      <label className="block mb-3">
        <span className="block text-sm text-ink-secondary mb-1">学習日</span>
        <input
          type="date"
          className="w-full border border-hairline rounded-xs px-2 py-1.5"
          {...register('studied_on', { required: '学習日を入力してください' })}
        />
        {errors.studied_on && <p className="text-error text-sm mt-1">{errors.studied_on.message}</p>}
      </label>
      <label className="block mb-3">
        <span className="block text-sm text-ink-secondary mb-1">学習時間(分)</span>
        <input
          type="number"
          className="w-full border border-hairline rounded-xs px-2 py-1.5"
          {...register('duration_min', {
            required: '学習時間を入力してください',
            valueAsNumber: true,
            min: { value: 1, message: '1分以上を入力してください' },
          })}
        />
        {errors.duration_min && <p className="text-error text-sm mt-1">{errors.duration_min.message}</p>}
      </label>
      <label className="block mb-3">
        <span className="block text-sm text-ink-secondary mb-1">メモ(任意)</span>
        <textarea className="w-full border border-hairline rounded-xs px-2 py-1.5" {...register('memo')} />
      </label>
      {serverError && <p className="text-error text-sm mb-3">{serverError}</p>}
      <button
        type="submit"
        disabled={isSubmitting}
        className="bg-primary text-white rounded-md px-4 py-2 font-medium disabled:opacity-50"
      >
        {submitLabel}
      </button>
    </form>
  )
}

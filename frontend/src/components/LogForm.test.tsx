import { describe, it, expect, vi } from 'vitest'
import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { LogForm } from './LogForm'
import type { Category } from '../types/domain'

const categories: Category[] = [
  { id: 'c1', user_id: null, name: '英語', created_at: '2026-01-01T00:00:00Z' },
  { id: 'c2', user_id: null, name: '数学', created_at: '2026-01-01T00:00:00Z' },
]

describe('LogForm', () => {
  it('カテゴリ未選択で送信するとバリデーションエラーを表示する', async () => {
    const onSubmit = vi.fn()
    const user = userEvent.setup()

    render(<LogForm categories={categories} onSubmit={onSubmit} submitLabel="記録する" />)
    await user.click(screen.getByRole('button', { name: '記録する' }))

    expect(await screen.findByText('カテゴリを選択してください')).toBeInTheDocument()
    expect(onSubmit).not.toHaveBeenCalled()
  })

  it('フォーム入力して送信するとonSubmitに値を渡す', async () => {
    const onSubmit = vi.fn().mockResolvedValue(undefined)
    const user = userEvent.setup()

    render(<LogForm categories={categories} onSubmit={onSubmit} submitLabel="記録する" />)

    await user.selectOptions(screen.getByLabelText('カテゴリ'), 'c1')
    await user.clear(screen.getByLabelText('学習時間(分)'))
    await user.type(screen.getByLabelText('学習時間(分)'), '45')
    await user.click(screen.getByRole('button', { name: '記録する' }))

    expect(onSubmit).toHaveBeenCalledTimes(1)
    const submitted = onSubmit.mock.calls[0][0]
    expect(submitted.category_id).toBe('c1')
    expect(submitted.duration_min).toBe(45)
  })

  it('defaultValuesを渡すとフォームに初期値が入る', () => {
    render(
      <LogForm
        categories={categories}
        defaultValues={{ category_id: 'c2', studied_on: '2026-08-01', duration_min: 60, memo: '復習' }}
        onSubmit={vi.fn()}
        submitLabel="更新する"
      />
    )

    expect(screen.getByLabelText('カテゴリ')).toHaveValue('c2')
    expect(screen.getByLabelText('学習時間(分)')).toHaveValue(60)
    expect(screen.getByLabelText('メモ(任意)')).toHaveValue('復習')
  })
})

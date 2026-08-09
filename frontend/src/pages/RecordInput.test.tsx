import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { MemoryRouter } from 'react-router'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { RecordInput } from './RecordInput'
import * as categoriesApi from '../api/categories'
import * as logsApi from '../api/logs'

vi.mock('../api/categories')
vi.mock('../api/logs')

function renderWithProviders() {
  const queryClient = new QueryClient({ defaultOptions: { queries: { retry: false } } })
  return render(
    <QueryClientProvider client={queryClient}>
      <MemoryRouter>
        <RecordInput />
      </MemoryRouter>
    </QueryClientProvider>
  )
}

describe('RecordInput', () => {
  beforeEach(() => {
    vi.mocked(categoriesApi.listCategories).mockReset()
    vi.mocked(logsApi.createLog).mockReset()
  })

  it('カテゴリ読み込み中はLoadingStateを表示する', () => {
    vi.mocked(categoriesApi.listCategories).mockReturnValue(new Promise(() => {}))

    renderWithProviders()

    expect(screen.getByText('読み込み中...')).toBeInTheDocument()
  })

  it('カテゴリ取得後、フォームに送信するとcreateLogを呼ぶ', async () => {
    vi.mocked(categoriesApi.listCategories).mockResolvedValue([
      { id: 'c1', user_id: null, name: '英語', created_at: '2026-01-01T00:00:00Z' },
    ])
    vi.mocked(logsApi.createLog).mockResolvedValue({
      id: '1',
      user_id: 'u1',
      category_id: 'c1',
      studied_on: '2026-08-09T00:00:00Z',
      duration_min: 45,
      created_at: '2026-08-09T00:00:00Z',
      updated_at: '2026-08-09T00:00:00Z',
    })
    const user = userEvent.setup()

    renderWithProviders()

    await waitFor(() => expect(screen.getByLabelText('カテゴリ')).toBeInTheDocument())
    await user.selectOptions(screen.getByLabelText('カテゴリ'), 'c1')
    await user.click(screen.getByRole('button', { name: '記録する' }))

    await waitFor(() => expect(logsApi.createLog).toHaveBeenCalledTimes(1))
  })
})

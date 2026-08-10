import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { MemoryRouter } from 'react-router'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { LogHistory } from './LogHistory'
import * as categoriesApi from '../api/categories'
import * as logsApi from '../api/logs'

vi.mock('../api/categories')
vi.mock('../api/logs')

function renderWithProviders() {
  const queryClient = new QueryClient({ defaultOptions: { queries: { retry: false } } })
  return render(
    <QueryClientProvider client={queryClient}>
      <MemoryRouter>
        <LogHistory />
      </MemoryRouter>
    </QueryClientProvider>
  )
}

const sampleLog = {
  id: 'l1',
  user_id: 'u1',
  category_id: 'c1',
  studied_on: '2026-08-01T00:00:00Z',
  duration_min: 30,
  memo: 'テスト',
  created_at: '2026-08-01T00:00:00Z',
  updated_at: '2026-08-01T00:00:00Z',
}

describe('LogHistory', () => {
  beforeEach(() => {
    vi.mocked(categoriesApi.listCategories).mockResolvedValue([
      { id: 'c1', user_id: null, name: '英語', created_at: '2026-01-01T00:00:00Z' },
    ])
    vi.mocked(logsApi.listLogs).mockReset().mockResolvedValue([sampleLog])
    vi.mocked(logsApi.updateLog).mockReset()
    vi.mocked(logsApi.deleteLog).mockReset()
  })

  afterEach(() => {
    vi.restoreAllMocks()
  })

  it('ログ一覧を表示する', async () => {
    renderWithProviders()

    expect(await screen.findByText('30分')).toBeInTheDocument()
    expect(screen.getByText('英語')).toBeInTheDocument()
  })

  it('編集ボタンでモーダルが開き、送信するとupdateLogを呼ぶ', async () => {
    vi.mocked(logsApi.updateLog).mockResolvedValue(sampleLog)
    const user = userEvent.setup()

    renderWithProviders()
    await screen.findByText('30分')
    await user.click(screen.getByRole('button', { name: '編集' }))

    expect(await screen.findByText('記録を編集')).toBeInTheDocument()
    await user.click(screen.getByRole('button', { name: '更新する' }))

    await waitFor(() => expect(logsApi.updateLog).toHaveBeenCalledWith('l1', expect.any(Object)))
  })

  it('削除ボタンで確認ダイアログが表示され、OKするとdeleteLogを呼ぶ', async () => {
    vi.spyOn(window, 'confirm').mockReturnValue(true)
    vi.mocked(logsApi.deleteLog).mockResolvedValue(undefined)
    const user = userEvent.setup()

    renderWithProviders()
    await screen.findByText('30分')
    await user.click(screen.getByRole('button', { name: '削除' }))

    expect(window.confirm).toHaveBeenCalled()
    await waitFor(() => expect(logsApi.deleteLog).toHaveBeenCalledWith('l1'))
  })

  it('削除確認をキャンセルするとdeleteLogを呼ばない', async () => {
    vi.spyOn(window, 'confirm').mockReturnValue(false)
    const user = userEvent.setup()

    renderWithProviders()
    await screen.findByText('30分')
    await user.click(screen.getByRole('button', { name: '削除' }))

    expect(logsApi.deleteLog).not.toHaveBeenCalled()
  })
})

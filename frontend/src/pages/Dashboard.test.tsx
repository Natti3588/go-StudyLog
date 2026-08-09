import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { MemoryRouter } from 'react-router'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { Dashboard } from './Dashboard'
import * as statsApi from '../api/stats'
import * as goalsApi from '../api/goals'
import * as logsApi from '../api/logs'
import * as AuthContextModule from '../context/AuthContext'

vi.mock('../api/stats')
vi.mock('../api/goals')
vi.mock('../api/logs')
vi.mock('../context/AuthContext', async () => {
  const actual = await vi.importActual<typeof AuthContextModule>('../context/AuthContext')
  return { ...actual, useAuth: vi.fn() }
})

function renderWithProviders() {
  const queryClient = new QueryClient({ defaultOptions: { queries: { retry: false } } })
  return render(
    <QueryClientProvider client={queryClient}>
      <MemoryRouter>
        <Dashboard />
      </MemoryRouter>
    </QueryClientProvider>
  )
}

describe('Dashboard', () => {
  beforeEach(() => {
    vi.mocked(AuthContextModule.useAuth).mockReturnValue({
      status: 'authenticated',
      user: { id: '1', email: 'a@example.com', is_admin: false, created_at: '2026-01-01T00:00:00Z' },
      login: vi.fn(),
      logout: vi.fn(),
    })
    vi.mocked(statsApi.getStatsSummary).mockReset().mockResolvedValue({
      total_min: 500,
      current_streak: 4,
      longest_streak: 10,
      weekly_target_min: 300,
      weekly_actual_min: 150,
    })
    vi.mocked(logsApi.listLogs).mockReset().mockResolvedValue([
      {
        id: 'l1',
        user_id: 'u1',
        category_id: 'c1',
        studied_on: '2026-08-08T00:00:00Z',
        duration_min: 30,
        created_at: '2026-08-08T00:00:00Z',
        updated_at: '2026-08-08T00:00:00Z',
      },
    ])
    vi.mocked(goalsApi.setWeeklyGoal).mockReset()
  })

  afterEach(() => {
    vi.restoreAllMocks()
  })

  it('サマリーと週間目標の進捗、最近の記録を表示する', async () => {
    renderWithProviders()

    expect(await screen.findByText('500分')).toBeInTheDocument()
    expect(screen.getByText('週間目標: 150 / 300分')).toBeInTheDocument()
    expect(screen.getByText('30分')).toBeInTheDocument()
  })

  it('「目標を変更」から新しい目標を送信するとsetWeeklyGoalを呼ぶ', async () => {
    vi.mocked(goalsApi.setWeeklyGoal).mockResolvedValue({
      user_id: 'u1',
      week_start: '2026-08-03T00:00:00Z',
      target_min: 400,
    })
    const user = userEvent.setup()

    renderWithProviders()
    await screen.findByText('500分')
    await user.click(screen.getByText('目標を変更'))
    await user.clear(screen.getByRole('spinbutton'))
    await user.type(screen.getByRole('spinbutton'), '400')
    await user.click(screen.getByRole('button', { name: '保存' }))

    await waitFor(() =>
      expect(goalsApi.setWeeklyGoal).toHaveBeenCalledWith(
        expect.objectContaining({ target_min: 400 })
      )
    )
  })

  it('記録が0件のときは「まだ記録がありません」を表示する', async () => {
    vi.mocked(logsApi.listLogs).mockResolvedValue([])

    renderWithProviders()

    expect(await screen.findByText('まだ記録がありません')).toBeInTheDocument()
  })
})

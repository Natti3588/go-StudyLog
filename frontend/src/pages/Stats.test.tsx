import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { Stats } from './Stats'
import * as statsApi from '../api/stats'

vi.mock('../api/stats')

function renderWithProviders() {
  const queryClient = new QueryClient({ defaultOptions: { queries: { retry: false } } })
  return render(
    <QueryClientProvider client={queryClient}>
      <Stats />
    </QueryClientProvider>
  )
}

describe('Stats', () => {
  beforeEach(() => {
    vi.mocked(statsApi.getStatsSummary).mockReset().mockResolvedValue({
      total_min: 500,
      current_streak: 4,
      longest_streak: 10,
      weekly_target_min: 300,
      weekly_actual_min: 150,
    })
    vi.mocked(statsApi.getHeatmap).mockReset().mockResolvedValue([])
  })

  it('サマリー数値を表示する', async () => {
    renderWithProviders()

    expect(await screen.findByText('500分')).toBeInTheDocument()
    expect(screen.getByText('4日')).toBeInTheDocument()
    expect(screen.getByText('10日')).toBeInTheDocument()
  })

  it('年を変更するとgetHeatmapが新しい年で呼ばれる', async () => {
    const user = userEvent.setup()

    renderWithProviders()
    await screen.findByText('500分')

    const currentYear = new Date().getUTCFullYear()
    await user.selectOptions(screen.getByLabelText('年'), String(currentYear - 1))

    await waitFor(() => expect(statsApi.getHeatmap).toHaveBeenCalledWith(currentYear - 1))
  })
})

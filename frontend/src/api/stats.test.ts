import { describe, it, expect, vi, beforeEach } from 'vitest'
import * as client from './client'
import { getStatsSummary, getHeatmap } from './stats'

vi.mock('./client', () => ({
  apiRequest: vi.fn(),
}))

describe('stats api', () => {
  beforeEach(() => {
    vi.mocked(client.apiRequest).mockReset()
  })

  it('getStatsSummaryはGET /stats/summaryを叩く', async () => {
    vi.mocked(client.apiRequest).mockResolvedValue({
      total_min: 100,
      current_streak: 3,
      longest_streak: 5,
      weekly_target_min: 300,
      weekly_actual_min: 120,
    })

    await getStatsSummary()

    expect(client.apiRequest).toHaveBeenCalledWith('/stats/summary')
  })

  it('getHeatmapはGET /stats/heatmap?year=を叩く', async () => {
    vi.mocked(client.apiRequest).mockResolvedValue([])

    await getHeatmap(2026)

    expect(client.apiRequest).toHaveBeenCalledWith('/stats/heatmap?year=2026')
  })
})

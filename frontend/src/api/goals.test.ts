import { describe, it, expect, vi, beforeEach } from 'vitest'
import * as client from './client'
import { setWeeklyGoal } from './goals'

vi.mock('./client', () => ({
  apiRequest: vi.fn(),
}))

describe('goals api', () => {
  beforeEach(() => {
    vi.mocked(client.apiRequest).mockReset()
  })

  it('setWeeklyGoalはPUT /goals/weeklyを叩く', async () => {
    vi.mocked(client.apiRequest).mockResolvedValue({ user_id: 'u1', week_start: '2026-08-10T00:00:00Z', target_min: 300 })

    await setWeeklyGoal({ week_start: '2026-08-10T00:00:00Z', target_min: 300 })

    expect(client.apiRequest).toHaveBeenCalledWith('/goals/weekly', {
      method: 'PUT',
      body: { week_start: '2026-08-10T00:00:00Z', target_min: 300 },
    })
  })
})

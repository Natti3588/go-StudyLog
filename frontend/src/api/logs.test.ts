import { describe, it, expect, vi, beforeEach } from 'vitest'
import * as client from './client'
import { listLogs, createLog, updateLog, deleteLog } from './logs'

vi.mock('./client', () => ({
  apiRequest: vi.fn(),
}))

const sampleInput = {
  category_id: 'c1',
  studied_on: '2026-08-09',
  duration_min: 30,
  memo: 'メモ',
}

describe('logs api', () => {
  beforeEach(() => {
    vi.mocked(client.apiRequest).mockReset()
  })

  it('listLogsはGET /logsを叩く', async () => {
    vi.mocked(client.apiRequest).mockResolvedValue([])

    await listLogs()

    expect(client.apiRequest).toHaveBeenCalledWith('/logs')
  })

  it('createLogはPOST /logsを叩く(studied_onはRFC3339に変換)', async () => {
    vi.mocked(client.apiRequest).mockResolvedValue({ id: '1' })

    await createLog(sampleInput)

    expect(client.apiRequest).toHaveBeenCalledWith('/logs', {
      method: 'POST',
      body: { ...sampleInput, studied_on: '2026-08-09T00:00:00Z' },
    })
  })

  it('updateLogはPUT /logs/{id}を叩く(studied_onはRFC3339に変換)', async () => {
    vi.mocked(client.apiRequest).mockResolvedValue({ id: '1' })

    await updateLog('1', sampleInput)

    expect(client.apiRequest).toHaveBeenCalledWith('/logs/1', {
      method: 'PUT',
      body: { ...sampleInput, studied_on: '2026-08-09T00:00:00Z' },
    })
  })

  it('deleteLogはDELETE /logs/{id}を叩く', async () => {
    vi.mocked(client.apiRequest).mockResolvedValue(undefined)

    await deleteLog('1')

    expect(client.apiRequest).toHaveBeenCalledWith('/logs/1', { method: 'DELETE' })
  })
})

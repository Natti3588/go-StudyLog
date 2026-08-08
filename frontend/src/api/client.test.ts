import { describe, it, expect, vi, beforeEach } from 'vitest'
import { apiRequest, ApiError } from './client'

describe('apiRequest', () => {
  beforeEach(() => {
    vi.restoreAllMocks()
  })

  it('GETリクエストにcredentials includeを付与する', async () => {
    const fetchMock = vi.fn().mockResolvedValue(
      new Response(JSON.stringify({ ok: true }), { status: 200 })
    )
    vi.stubGlobal('fetch', fetchMock)

    await apiRequest('/categories')

    expect(fetchMock).toHaveBeenCalledWith(
      expect.stringContaining('/categories'),
      expect.objectContaining({ credentials: 'include', method: 'GET' })
    )
  })

  it('bodyを渡すとJSON化してPOSTする', async () => {
    const fetchMock = vi.fn().mockResolvedValue(
      new Response(JSON.stringify({ id: '1' }), { status: 201 })
    )
    vi.stubGlobal('fetch', fetchMock)

    await apiRequest('/logs', { method: 'POST', body: { duration_min: 30 } })

    expect(fetchMock).toHaveBeenCalledWith(
      expect.stringContaining('/logs'),
      expect.objectContaining({
        method: 'POST',
        body: JSON.stringify({ duration_min: 30 }),
        headers: { 'Content-Type': 'application/json' },
      })
    )
  })

  it('200番台以外のレスポンスでApiErrorを投げる', async () => {
    vi.stubGlobal(
      'fetch',
      vi.fn().mockResolvedValue(new Response('bad request', { status: 400 }))
    )

    await expect(apiRequest('/logs', { method: 'POST', body: {} })).rejects.toThrow(ApiError)
  })

  it('401レスポンスでstatus=401のApiErrorを投げる', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue(new Response('', { status: 401 })))

    await expect(apiRequest('/me')).rejects.toMatchObject({ status: 401 })
  })

  it('204レスポンスはundefinedを返す', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue(new Response(null, { status: 204 })))

    const result = await apiRequest('/logs/1', { method: 'DELETE' })

    expect(result).toBeUndefined()
  })
})

import { describe, it, expect, vi, beforeEach } from 'vitest'
import * as client from './client'
import { listCategories } from './categories'

vi.mock('./client', () => ({
  apiRequest: vi.fn(),
}))

describe('categories api', () => {
  beforeEach(() => {
    vi.mocked(client.apiRequest).mockReset()
  })

  it('listCategoriesはGET /categoriesを叩く', async () => {
    vi.mocked(client.apiRequest).mockResolvedValue([{ id: 'c1', user_id: null, name: '英語', created_at: '2026-01-01T00:00:00Z' }])

    const result = await listCategories()

    expect(client.apiRequest).toHaveBeenCalledWith('/categories')
    expect(result).toEqual([{ id: 'c1', user_id: null, name: '英語', created_at: '2026-01-01T00:00:00Z' }])
  })
})

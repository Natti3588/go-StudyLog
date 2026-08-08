import { describe, it, expect, vi, beforeEach } from 'vitest'
import * as client from './client'
import { signup, login, logout, me } from './auth'

vi.mock('./client', () => ({
  apiRequest: vi.fn(),
}))

describe('auth api', () => {
  beforeEach(() => {
    vi.mocked(client.apiRequest).mockReset()
  })

  it('signupはPOST /signupを叩く', async () => {
    vi.mocked(client.apiRequest).mockResolvedValue({ id: '1' })

    await signup('a@example.com', 'password1')

    expect(client.apiRequest).toHaveBeenCalledWith('/signup', {
      method: 'POST',
      body: { email: 'a@example.com', password: 'password1' },
    })
  })

  it('loginはPOST /loginを叩く', async () => {
    vi.mocked(client.apiRequest).mockResolvedValue({ token: 't' })

    await login('a@example.com', 'password1')

    expect(client.apiRequest).toHaveBeenCalledWith('/login', {
      method: 'POST',
      body: { email: 'a@example.com', password: 'password1' },
    })
  })

  it('logoutはPOST /logoutを叩く', async () => {
    vi.mocked(client.apiRequest).mockResolvedValue(undefined)

    await logout()

    expect(client.apiRequest).toHaveBeenCalledWith('/logout', { method: 'POST' })
  })

  it('meはGET /meを叩く', async () => {
    vi.mocked(client.apiRequest).mockResolvedValue({ id: '1' })

    await me()

    expect(client.apiRequest).toHaveBeenCalledWith('/me')
  })
})

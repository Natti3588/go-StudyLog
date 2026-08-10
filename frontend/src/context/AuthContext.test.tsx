import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { AuthProvider, useAuth } from './AuthContext'
import * as authApi from '../api/auth'
import { ApiError, setUnauthorizedHandler } from '../api/client'

vi.mock('../api/auth')
vi.mock('../api/client', async () => {
  const actual = await vi.importActual<typeof import('../api/client')>('../api/client')
  return { ...actual, setUnauthorizedHandler: vi.fn() }
})

function Consumer() {
  const { status, user, login } = useAuth()
  return (
    <div>
      <span data-testid="status">{status}</span>
      <span data-testid="email">{user?.email ?? ''}</span>
      <button onClick={() => login('a@example.com', 'password1')}>login</button>
    </div>
  )
}

describe('AuthProvider', () => {
  beforeEach(() => {
    vi.mocked(authApi.me).mockReset()
    vi.mocked(authApi.login).mockReset()
    vi.mocked(setUnauthorizedHandler).mockClear()
  })

  it('起動時にGET /meが成功したらauthenticatedになる', async () => {
    vi.mocked(authApi.me).mockResolvedValue({
      id: '1',
      email: 'a@example.com',
      is_admin: false,
      created_at: '2026-01-01T00:00:00Z',
    })

    render(
      <AuthProvider>
        <Consumer />
      </AuthProvider>
    )

    expect(screen.getByTestId('status').textContent).toBe('loading')
    await waitFor(() => expect(screen.getByTestId('status').textContent).toBe('authenticated'))
    expect(screen.getByTestId('email').textContent).toBe('a@example.com')
  })

  it('起動時にGET /meが401ならunauthenticatedになる', async () => {
    vi.mocked(authApi.me).mockRejectedValue(new ApiError(401, 'unauthorized'))

    render(
      <AuthProvider>
        <Consumer />
      </AuthProvider>
    )

    await waitFor(() => expect(screen.getByTestId('status').textContent).toBe('unauthenticated'))
  })

  it('loginを呼ぶとauthenticatedに遷移する', async () => {
    vi.mocked(authApi.me)
      .mockRejectedValueOnce(new ApiError(401, 'unauthorized'))
      .mockResolvedValueOnce({
        id: '1',
        email: 'a@example.com',
        is_admin: false,
        created_at: '2026-01-01T00:00:00Z',
      })
    vi.mocked(authApi.login).mockResolvedValue({ token: 't' })

    const user = userEvent.setup()
    render(
      <AuthProvider>
        <Consumer />
      </AuthProvider>
    )

    await waitFor(() => expect(screen.getByTestId('status').textContent).toBe('unauthenticated'))
    await user.click(screen.getByText('login'))
    await waitFor(() => expect(screen.getByTestId('status').textContent).toBe('authenticated'))
  })

  it('マウント時にunauthorizedHandlerを登録する', async () => {
    vi.mocked(authApi.me).mockResolvedValue({
      id: '1',
      email: 'a@example.com',
      is_admin: false,
      created_at: '2026-01-01T00:00:00Z',
    })

    render(
      <AuthProvider>
        <Consumer />
      </AuthProvider>
    )

    await waitFor(() => expect(screen.getByTestId('status').textContent).toBe('authenticated'))
    expect(vi.mocked(setUnauthorizedHandler)).toHaveBeenCalledWith(expect.any(Function))
  })

  it('登録したunauthorizedHandlerを呼ぶとunauthenticatedになる', async () => {
    vi.mocked(authApi.me).mockResolvedValue({
      id: '1',
      email: 'a@example.com',
      is_admin: false,
      created_at: '2026-01-01T00:00:00Z',
    })

    render(
      <AuthProvider>
        <Consumer />
      </AuthProvider>
    )

    await waitFor(() => expect(screen.getByTestId('status').textContent).toBe('authenticated'))

    const registeredHandler = vi.mocked(setUnauthorizedHandler).mock.calls[0][0]
    registeredHandler?.()

    await waitFor(() => expect(screen.getByTestId('status').textContent).toBe('unauthenticated'))
  })
})

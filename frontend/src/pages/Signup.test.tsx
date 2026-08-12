import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { MemoryRouter } from 'react-router'
import { Signup } from './Signup'
import * as authApi from '../api/auth'
import * as AuthContextModule from '../context/AuthContext'
import { ApiError } from '../api/client'

const navigateMock = vi.fn()

vi.mock('../api/auth')

vi.mock('react-router', async () => {
  const actual = await vi.importActual<typeof import('react-router')>('react-router')
  return { ...actual, useNavigate: () => navigateMock }
})

vi.mock('../context/AuthContext', async () => {
  const actual = await vi.importActual<typeof AuthContextModule>('../context/AuthContext')
  return { ...actual, useAuth: vi.fn() }
})

function mockAuth(login = vi.fn()) {
  vi.mocked(AuthContextModule.useAuth).mockReturnValue({
    status: 'unauthenticated',
    user: null,
    login,
    logout: vi.fn(),
  })
  return login
}

describe('Signup', () => {
  beforeEach(() => {
    vi.mocked(authApi.signup).mockReset()
    vi.mocked(AuthContextModule.useAuth).mockReset()
    navigateMock.mockReset()
  })

  it('未入力で送信するとバリデーションエラーを表示する', async () => {
    mockAuth()
    const user = userEvent.setup()

    render(
      <MemoryRouter>
        <Signup />
      </MemoryRouter>
    )
    await user.click(screen.getByRole('button', { name: 'サインアップ' }))

    expect(await screen.findByText('メールアドレスを入力してください')).toBeInTheDocument()
  })

  it('メール重複エラー時に専用メッセージを表示する', async () => {
    vi.mocked(authApi.signup).mockRejectedValue(new ApiError(409, 'email already exists'))
    mockAuth()
    const user = userEvent.setup()

    render(
      <MemoryRouter>
        <Signup />
      </MemoryRouter>
    )
    await user.type(screen.getByLabelText('メールアドレス'), 'a@example.com')
    await user.type(screen.getByLabelText('パスワード'), 'password1')
    await user.click(screen.getByRole('button', { name: 'サインアップ' }))

    expect(
      await screen.findByText('このメールアドレスは既に登録されています')
    ).toBeInTheDocument()
  })

  it('登録に成功すると自動ログインしてダッシュボードへ遷移する', async () => {
    vi.mocked(authApi.signup).mockResolvedValue({
      id: 'u1',
      email: 'a@example.com',
      is_admin: false,
      created_at: '2026-08-12T00:00:00Z',
    })
    const loginMock = mockAuth(vi.fn().mockResolvedValue(undefined))
    const user = userEvent.setup()

    render(
      <MemoryRouter>
        <Signup />
      </MemoryRouter>
    )
    await user.type(screen.getByLabelText('メールアドレス'), 'a@example.com')
    await user.type(screen.getByLabelText('パスワード'), 'password1')
    await user.click(screen.getByRole('button', { name: 'サインアップ' }))

    await waitFor(() => expect(loginMock).toHaveBeenCalledWith('a@example.com', 'password1'))
    expect(navigateMock).toHaveBeenCalledWith('/', { replace: true })
  })

  it('登録に成功しても自動ログインに失敗した場合はログイン画面へ遷移する', async () => {
    vi.mocked(authApi.signup).mockResolvedValue({
      id: 'u1',
      email: 'a@example.com',
      is_admin: false,
      created_at: '2026-08-12T00:00:00Z',
    })
    mockAuth(vi.fn().mockRejectedValue(new ApiError(500, 'boom')))
    const user = userEvent.setup()

    render(
      <MemoryRouter>
        <Signup />
      </MemoryRouter>
    )
    await user.type(screen.getByLabelText('メールアドレス'), 'a@example.com')
    await user.type(screen.getByLabelText('パスワード'), 'password1')
    await user.click(screen.getByRole('button', { name: 'サインアップ' }))

    await waitFor(() => expect(navigateMock).toHaveBeenCalledWith('/login', { replace: true }))
  })
})

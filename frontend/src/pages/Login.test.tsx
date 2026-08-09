import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { MemoryRouter } from 'react-router'
import { Login } from './Login'
import * as AuthContextModule from '../context/AuthContext'
import { ApiError } from '../api/client'

vi.mock('../context/AuthContext', async () => {
  const actual = await vi.importActual<typeof AuthContextModule>('../context/AuthContext')
  return { ...actual, useAuth: vi.fn() }
})

describe('Login', () => {
  beforeEach(() => {
    vi.mocked(AuthContextModule.useAuth).mockReset()
  })

  it('未入力で送信するとバリデーションエラーを表示する', async () => {
    vi.mocked(AuthContextModule.useAuth).mockReturnValue({
      status: 'unauthenticated',
      user: null,
      login: vi.fn(),
      logout: vi.fn(),
    })
    const user = userEvent.setup()

    render(
      <MemoryRouter>
        <Login />
      </MemoryRouter>
    )
    await user.click(screen.getByRole('button', { name: 'ログイン' }))

    expect(await screen.findByText('メールアドレスを入力してください')).toBeInTheDocument()
  })

  it('ログイン失敗時にサーバーエラーを表示する', async () => {
    const loginMock = vi.fn().mockRejectedValue(new ApiError(401, 'unauthorized'))
    vi.mocked(AuthContextModule.useAuth).mockReturnValue({
      status: 'unauthenticated',
      user: null,
      login: loginMock,
      logout: vi.fn(),
    })
    const user = userEvent.setup()

    render(
      <MemoryRouter>
        <Login />
      </MemoryRouter>
    )
    await user.type(screen.getByLabelText('メールアドレス'), 'a@example.com')
    await user.type(screen.getByLabelText('パスワード'), 'wrongpass')
    await user.click(screen.getByRole('button', { name: 'ログイン' }))

    expect(
      await screen.findByText('メールアドレスまたはパスワードが正しくありません')
    ).toBeInTheDocument()
  })
})

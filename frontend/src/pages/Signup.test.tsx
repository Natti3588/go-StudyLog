import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { MemoryRouter } from 'react-router'
import { Signup } from './Signup'
import * as authApi from '../api/auth'
import { ApiError } from '../api/client'

vi.mock('../api/auth')

describe('Signup', () => {
  beforeEach(() => {
    vi.mocked(authApi.signup).mockReset()
  })

  it('未入力で送信するとバリデーションエラーを表示する', async () => {
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
})

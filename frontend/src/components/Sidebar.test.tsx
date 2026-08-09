import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen } from '@testing-library/react'
import { MemoryRouter } from 'react-router'
import { Sidebar } from './Sidebar'
import * as AuthContextModule from '../context/AuthContext'

vi.mock('../context/AuthContext', async () => {
  const actual = await vi.importActual<typeof AuthContextModule>('../context/AuthContext')
  return { ...actual, useAuth: vi.fn() }
})

describe('Sidebar', () => {
  beforeEach(() => {
    vi.mocked(AuthContextModule.useAuth).mockReturnValue({
      status: 'authenticated',
      user: { id: '1', email: 'a@example.com', is_admin: false, created_at: '2026-01-01T00:00:00Z' },
      login: vi.fn(),
      logout: vi.fn(),
    })
  })

  it('4つのナビリンクとログアウトボタンを表示する', () => {
    render(
      <MemoryRouter initialEntries={['/']}>
        <Sidebar />
      </MemoryRouter>
    )

    expect(screen.getByText('ダッシュボード')).toBeInTheDocument()
    expect(screen.getByText('記録入力')).toBeInTheDocument()
    expect(screen.getByText('ログ履歴')).toBeInTheDocument()
    expect(screen.getByText('統計')).toBeInTheDocument()
    expect(screen.getByText('ログアウト')).toBeInTheDocument()
  })

  it('現在のパスに対応するリンクにaria-current="page"が付く', () => {
    render(
      <MemoryRouter initialEntries={['/logs']}>
        <Sidebar />
      </MemoryRouter>
    )

    expect(screen.getByText('ログ履歴')).toHaveAttribute('aria-current', 'page')
    expect(screen.getByText('ダッシュボード')).not.toHaveAttribute('aria-current')
  })
})

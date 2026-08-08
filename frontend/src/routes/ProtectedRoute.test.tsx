import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen } from '@testing-library/react'
import { MemoryRouter, Routes, Route } from 'react-router-dom'
import { ProtectedRoute } from './ProtectedRoute'
import * as AuthContextModule from '../context/AuthContext'

vi.mock('../context/AuthContext', async () => {
  const actual = await vi.importActual<typeof AuthContextModule>('../context/AuthContext')
  return { ...actual, useAuth: vi.fn() }
})

function renderWithRouter(initialPath: string) {
  return render(
    <MemoryRouter initialEntries={[initialPath]}>
      <Routes>
        <Route
          path="/"
          element={
            <ProtectedRoute>
              <div>保護されたコンテンツ</div>
            </ProtectedRoute>
          }
        />
        <Route path="/login" element={<div>ログイン画面</div>} />
      </Routes>
    </MemoryRouter>
  )
}

describe('ProtectedRoute', () => {
  beforeEach(() => {
    vi.mocked(AuthContextModule.useAuth).mockReset()
  })

  it('loading中はローディング表示を出す', () => {
    vi.mocked(AuthContextModule.useAuth).mockReturnValue({
      status: 'loading',
      user: null,
      login: vi.fn(),
      logout: vi.fn(),
    })

    renderWithRouter('/')

    expect(screen.getByText('読み込み中...')).toBeInTheDocument()
  })

  it('unauthenticatedなら/loginにリダイレクトする', () => {
    vi.mocked(AuthContextModule.useAuth).mockReturnValue({
      status: 'unauthenticated',
      user: null,
      login: vi.fn(),
      logout: vi.fn(),
    })

    renderWithRouter('/')

    expect(screen.getByText('ログイン画面')).toBeInTheDocument()
  })

  it('authenticatedならchildrenを表示する', () => {
    vi.mocked(AuthContextModule.useAuth).mockReturnValue({
      status: 'authenticated',
      user: { id: '1', email: 'a@example.com', is_admin: false, created_at: '2026-01-01T00:00:00Z' },
      login: vi.fn(),
      logout: vi.fn(),
    })

    renderWithRouter('/')

    expect(screen.getByText('保護されたコンテンツ')).toBeInTheDocument()
  })
})

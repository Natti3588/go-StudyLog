import { describe, it, expect } from 'vitest'
import { render, screen } from '@testing-library/react'
import { MemoryRouter } from 'react-router'
import { BottomNav } from './BottomNav'

describe('BottomNav', () => {
  it('4つのナビリンクを表示する', () => {
    render(
      <MemoryRouter initialEntries={['/']}>
        <BottomNav />
      </MemoryRouter>
    )

    expect(screen.getByText('ホーム')).toBeInTheDocument()
    expect(screen.getByText('記録')).toBeInTheDocument()
    expect(screen.getByText('履歴')).toBeInTheDocument()
    expect(screen.getByText('統計')).toBeInTheDocument()
  })

  it('現在のパスに対応するリンクにaria-current="page"が付く', () => {
    render(
      <MemoryRouter initialEntries={['/stats']}>
        <BottomNav />
      </MemoryRouter>
    )

    expect(screen.getByText('統計')).toHaveAttribute('aria-current', 'page')
  })

  it('/logs/newでは記録のみがアクティブになり履歴はアクティブにならない', () => {
    render(
      <MemoryRouter initialEntries={['/logs/new']}>
        <BottomNav />
      </MemoryRouter>
    )

    expect(screen.getByText('記録')).toHaveAttribute('aria-current', 'page')
    expect(screen.getByText('履歴')).not.toHaveAttribute('aria-current')
  })
})

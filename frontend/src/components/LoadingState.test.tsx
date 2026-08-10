import { describe, it, expect } from 'vitest'
import { render, screen } from '@testing-library/react'
import { LoadingState } from './LoadingState'

describe('LoadingState', () => {
  it('読み込み中メッセージを表示する', () => {
    render(<LoadingState />)

    expect(screen.getByText('読み込み中...')).toBeInTheDocument()
  })
})

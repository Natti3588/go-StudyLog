import { describe, it, expect, vi } from 'vitest'
import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { ErrorState } from './ErrorState'

describe('ErrorState', () => {
  it('エラーメッセージと再読み込みボタンを表示し、クリックでonRetryを呼ぶ', async () => {
    const onRetry = vi.fn()
    const user = userEvent.setup()

    render(<ErrorState onRetry={onRetry} />)

    expect(screen.getByText('エラーが発生しました')).toBeInTheDocument()
    await user.click(screen.getByRole('button', { name: '再読み込み' }))

    expect(onRetry).toHaveBeenCalledTimes(1)
  })
})

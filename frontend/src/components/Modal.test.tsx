import { describe, it, expect, vi } from 'vitest'
import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { Modal } from './Modal'

describe('Modal', () => {
  it('childrenを表示する', () => {
    render(
      <Modal onClose={vi.fn()}>
        <p>モーダルの中身</p>
      </Modal>
    )

    expect(screen.getByText('モーダルの中身')).toBeInTheDocument()
  })

  it('オーバーレイをクリックするとonCloseを呼ぶ', async () => {
    const onClose = vi.fn()
    const user = userEvent.setup()

    render(
      <Modal onClose={onClose}>
        <p>モーダルの中身</p>
      </Modal>
    )
    await user.click(screen.getByTestId('modal-overlay'))

    expect(onClose).toHaveBeenCalledTimes(1)
  })

  it('中身をクリックしてもonCloseを呼ばない', async () => {
    const onClose = vi.fn()
    const user = userEvent.setup()

    render(
      <Modal onClose={onClose}>
        <p>モーダルの中身</p>
      </Modal>
    )
    await user.click(screen.getByText('モーダルの中身'))

    expect(onClose).not.toHaveBeenCalled()
  })
})

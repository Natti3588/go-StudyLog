import type { ReactNode } from 'react'

interface ModalProps {
  onClose: () => void
  children: ReactNode
}

export function Modal({ onClose, children }: ModalProps) {
  return (
    <div
      data-testid="modal-overlay"
      className="fixed inset-0 bg-black/30 flex items-center justify-center p-4"
      onClick={onClose}
    >
      <div className="bg-surface rounded-xl p-6 w-full max-w-md" onClick={(e) => e.stopPropagation()}>
        {children}
      </div>
    </div>
  )
}

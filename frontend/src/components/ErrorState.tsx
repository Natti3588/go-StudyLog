interface ErrorStateProps {
  onRetry: () => void
}

export function ErrorState({ onRetry }: ErrorStateProps) {
  return (
    <div className="flex flex-col items-center justify-center p-12 gap-3">
      <p className="text-error">エラーが発生しました</p>
      <button
        onClick={onRetry}
        className="bg-surface border border-hairline rounded-md px-4 py-2 text-ink"
      >
        再読み込み
      </button>
    </div>
  )
}

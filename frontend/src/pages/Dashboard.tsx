import { useAuth } from '../context/AuthContext'

export function Dashboard() {
  const { user, logout } = useAuth()

  return (
    <div className="p-6">
      <h1 className="text-xl font-bold text-ink">ようこそ、{user?.email}</h1>
      <button
        onClick={() => logout()}
        className="mt-4 bg-surface border border-hairline rounded-md px-4 py-2 text-ink"
      >
        ログアウト
      </button>
    </div>
  )
}

import { NavLink } from 'react-router'
import { useAuth } from '../context/AuthContext'

const navItems = [
  { to: '/', label: 'ダッシュボード', end: true },
  { to: '/logs/new', label: '記録入力', end: false },
  { to: '/logs', label: 'ログ履歴', end: false },
  { to: '/stats', label: '統計', end: false },
]

export function Sidebar() {
  const { logout } = useAuth()

  return (
    <aside className="hidden md:flex md:flex-col md:w-60 md:shrink-0 bg-canvas border-r border-hairline p-4">
      <span className="text-lg font-bold text-ink mb-6">StudyLog</span>
      <nav className="flex flex-col gap-1">
        {navItems.map((item) => (
          <NavLink
            key={item.to}
            to={item.to}
            end={item.end}
            className={({ isActive }) =>
              `rounded-md px-3 py-2 text-sm ${
                isActive ? 'bg-canvas-soft text-primary font-medium' : 'text-ink'
              }`
            }
          >
            {item.label}
          </NavLink>
        ))}
      </nav>
      <button onClick={() => logout()} className="mt-auto text-sm text-ink-muted text-left px-3 py-2">
        ログアウト
      </button>
    </aside>
  )
}

import { NavLink } from 'react-router'

const navItems = [
  { to: '/', label: 'ホーム', end: true },
  { to: '/logs/new', label: '記録', end: false },
  { to: '/logs', label: '履歴', end: true },
  { to: '/stats', label: '統計', end: false },
]

export function BottomNav() {
  return (
    <nav className="md:hidden fixed bottom-0 left-0 right-0 flex bg-canvas border-t border-hairline">
      {navItems.map((item) => (
        <NavLink
          key={item.to}
          to={item.to}
          end={item.end}
          className={({ isActive }) =>
            `flex-1 text-center py-2 text-xs ${isActive ? 'text-primary font-medium' : 'text-ink-muted'}`
          }
        >
          {item.label}
        </NavLink>
      ))}
    </nav>
  )
}

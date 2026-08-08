import { Outlet } from 'react-router'

export function AppLayout() {
  return (
    <div className="min-h-screen bg-canvas-soft">
      <Outlet />
    </div>
  )
}

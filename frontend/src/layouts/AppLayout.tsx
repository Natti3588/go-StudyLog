import { Outlet } from 'react-router'
import { Sidebar } from '../components/Sidebar'
import { BottomNav } from '../components/BottomNav'

export function AppLayout() {
  return (
    <div className="min-h-screen bg-canvas-soft md:flex">
      <Sidebar />
      <main className="flex-1 pb-16 md:pb-0">
        <Outlet />
      </main>
      <BottomNav />
    </div>
  )
}

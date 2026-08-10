import { createBrowserRouter } from 'react-router'
import { AppLayout } from '../layouts/AppLayout'
import { AuthLayout } from '../layouts/AuthLayout'
import { ProtectedRoute } from './ProtectedRoute'
import { Login } from '../pages/Login'
import { Signup } from '../pages/Signup'
import { Dashboard } from '../pages/Dashboard'
import { RecordInput } from '../pages/RecordInput'
import { LogHistory } from '../pages/LogHistory'
import { Stats } from '../pages/Stats'

export const router = createBrowserRouter([
  {
    element: <AuthLayout />,
    children: [
      { path: '/login', element: <Login /> },
      { path: '/signup', element: <Signup /> },
    ],
  },
  {
    element: (
      <ProtectedRoute>
        <AppLayout />
      </ProtectedRoute>
    ),
    children: [
      { path: '/', element: <Dashboard /> },
      { path: '/logs/new', element: <RecordInput /> },
      { path: '/logs', element: <LogHistory /> },
      { path: '/stats', element: <Stats /> },
    ],
  },
])

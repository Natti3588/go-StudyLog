import { createBrowserRouter } from 'react-router'
import { AppLayout } from '../layouts/AppLayout'
import { AuthLayout } from '../layouts/AuthLayout'
import { ProtectedRoute } from './ProtectedRoute'
import { LoadingState } from '../components/LoadingState'

// 各ページはルートに一致したときだけ読み込む(初回ロードのJSを小さく保つため)
export const router = createBrowserRouter([
  {
    element: <AuthLayout />,
    HydrateFallback: LoadingState,
    children: [
      {
        path: '/login',
        lazy: async () => ({ Component: (await import('../pages/Login')).Login }),
      },
      {
        path: '/signup',
        lazy: async () => ({ Component: (await import('../pages/Signup')).Signup }),
      },
    ],
  },
  {
    element: (
      <ProtectedRoute>
        <AppLayout />
      </ProtectedRoute>
    ),
    HydrateFallback: LoadingState,
    children: [
      {
        path: '/',
        lazy: async () => ({ Component: (await import('../pages/Dashboard')).Dashboard }),
      },
      {
        path: '/logs/new',
        lazy: async () => ({ Component: (await import('../pages/RecordInput')).RecordInput }),
      },
      {
        path: '/logs',
        lazy: async () => ({ Component: (await import('../pages/LogHistory')).LogHistory }),
      },
      {
        path: '/stats',
        lazy: async () => ({ Component: (await import('../pages/Stats')).Stats }),
      },
    ],
  },
])

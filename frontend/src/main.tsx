import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import { RouterProvider } from 'react-router/dom'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { AuthProvider } from './context/AuthContext'
import { router } from './routes/router'
import './index.css'

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      // 既定は0なので画面を移動するたびに全クエリが再フェッチされる。
      // 学習記録はリアルタイム性が不要なため5分間はキャッシュを使い回す。
      staleTime: 5 * 60 * 1000,
      refetchOnWindowFocus: false,
      retry: 1,
    },
  },
})

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <QueryClientProvider client={queryClient}>
      <AuthProvider>
        <RouterProvider router={router} />
      </AuthProvider>
    </QueryClientProvider>
  </StrictMode>
)

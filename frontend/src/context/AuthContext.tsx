import { createContext, useContext, useEffect, useState, type ReactNode } from 'react'
import type { User } from '../types/domain'
import { me as fetchMe, login as apiLogin, logout as apiLogout } from '../api/auth'
import { setUnauthorizedHandler } from '../api/client'

type AuthStatus = 'loading' | 'authenticated' | 'unauthenticated'

interface AuthContextValue {
  status: AuthStatus
  user: User | null
  login: (email: string, password: string) => Promise<void>
  logout: () => Promise<void>
}

const AuthContext = createContext<AuthContextValue | undefined>(undefined)

export function AuthProvider({ children }: { children: ReactNode }) {
  const [status, setStatus] = useState<AuthStatus>('loading')
  const [user, setUser] = useState<User | null>(null)

  useEffect(() => {
    setUnauthorizedHandler(() => {
      setUser(null)
      setStatus('unauthenticated')
    })

    fetchMe()
      .then((u) => {
        setUser(u)
        setStatus('authenticated')
      })
      .catch(() => {
        setStatus('unauthenticated')
      })

    return () => setUnauthorizedHandler(null)
  }, [])

  async function login(email: string, password: string) {
    await apiLogin(email, password)
    const u = await fetchMe()
    setUser(u)
    setStatus('authenticated')
  }

  async function logout() {
    await apiLogout()
    setUser(null)
    setStatus('unauthenticated')
  }

  return (
    <AuthContext.Provider value={{ status, user, login, logout }}>
      {children}
    </AuthContext.Provider>
  )
}

export function useAuth(): AuthContextValue {
  const ctx = useContext(AuthContext)
  if (!ctx) {
    throw new Error('useAuth must be used within AuthProvider')
  }
  return ctx
}

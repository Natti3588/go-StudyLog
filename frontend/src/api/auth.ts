import { apiRequest } from './client'
import type { User } from '../types/domain'

export function signup(email: string, password: string): Promise<User> {
  return apiRequest<User>('/signup', { method: 'POST', body: { email, password } })
}

export function login(email: string, password: string): Promise<{ token: string }> {
  return apiRequest<{ token: string }>('/login', { method: 'POST', body: { email, password } })
}

export function logout(): Promise<void> {
  return apiRequest<void>('/logout', { method: 'POST' })
}

export function me(): Promise<User> {
  return apiRequest<User>('/me')
}

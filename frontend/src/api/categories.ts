import { apiRequest } from './client'
import type { Category } from '../types/domain'

export function listCategories(): Promise<Category[]> {
  return apiRequest<Category[]>('/categories')
}

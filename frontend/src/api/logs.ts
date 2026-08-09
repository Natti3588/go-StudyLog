import { apiRequest } from './client'
import type { StudyLog } from '../types/domain'

export interface LogInput {
  category_id: string
  studied_on: string
  duration_min: number
  memo?: string
}

export function listLogs(): Promise<StudyLog[]> {
  return apiRequest<StudyLog[]>('/logs')
}

export function createLog(input: LogInput): Promise<StudyLog> {
  return apiRequest<StudyLog>('/logs', { method: 'POST', body: input })
}

export function updateLog(id: string, input: LogInput): Promise<StudyLog> {
  return apiRequest<StudyLog>(`/logs/${id}`, { method: 'PUT', body: input })
}

export function deleteLog(id: string): Promise<void> {
  return apiRequest<void>(`/logs/${id}`, { method: 'DELETE' })
}

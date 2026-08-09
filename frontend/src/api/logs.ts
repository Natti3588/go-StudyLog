import { apiRequest } from './client'
import type { StudyLog } from '../types/domain'

export interface LogInput {
  category_id: string
  studied_on: string
  duration_min: number
  memo?: string
}

function toRFC3339(dateStr: string): string {
  return `${dateStr}T00:00:00Z`
}

export function listLogs(): Promise<StudyLog[]> {
  return apiRequest<StudyLog[]>('/logs')
}

export function createLog(input: LogInput): Promise<StudyLog> {
  return apiRequest<StudyLog>('/logs', {
    method: 'POST',
    body: { ...input, studied_on: toRFC3339(input.studied_on) },
  })
}

export function updateLog(id: string, input: LogInput): Promise<StudyLog> {
  return apiRequest<StudyLog>(`/logs/${id}`, {
    method: 'PUT',
    body: { ...input, studied_on: toRFC3339(input.studied_on) },
  })
}

export function deleteLog(id: string): Promise<void> {
  return apiRequest<void>(`/logs/${id}`, { method: 'DELETE' })
}

import { describe, it, expect } from 'vitest'
import { localDateString } from './date'

describe('localDateString', () => {
  it('通常の日付をYYYY-MM-DD形式で返す', () => {
    expect(localDateString(new Date(2026, 7, 9))).toBe('2026-08-09')
  })

  it('月・日が1桁の場合は0埋めする', () => {
    expect(localDateString(new Date(2026, 0, 5))).toBe('2026-01-05')
  })

  it('ローカルタイムで深夜0時台でもUTC変換せずローカルの年月日を返す', () => {
    expect(localDateString(new Date(2026, 0, 5, 0, 30))).toBe('2026-01-05')
  })
})

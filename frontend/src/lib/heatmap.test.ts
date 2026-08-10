import { describe, it, expect } from 'vitest'
import { levelForMinutes, buildHeatmapData } from './heatmap'

describe('levelForMinutes', () => {
  it('0分はlevel 0', () => {
    expect(levelForMinutes(0)).toBe(0)
  })
  it('1〜29分はlevel 1', () => {
    expect(levelForMinutes(1)).toBe(1)
    expect(levelForMinutes(29)).toBe(1)
  })
  it('30〜59分はlevel 2', () => {
    expect(levelForMinutes(30)).toBe(2)
    expect(levelForMinutes(59)).toBe(2)
  })
  it('60〜119分はlevel 3', () => {
    expect(levelForMinutes(60)).toBe(3)
    expect(levelForMinutes(119)).toBe(3)
  })
  it('120分以上はlevel 4', () => {
    expect(levelForMinutes(120)).toBe(4)
    expect(levelForMinutes(500)).toBe(4)
  })
})

describe('buildHeatmapData', () => {
  it('指定年の1/1から12/31まで全日程を生成する(平年365日)', () => {
    const result = buildHeatmapData(2026, [])

    expect(result).toHaveLength(365)
    expect(result[0]).toEqual({ date: '2026-01-01', count: 0, level: 0 })
    expect(result[364]).toEqual({ date: '2026-12-31', count: 0, level: 0 })
  })

  it('バックエンドのデータを対応する日付にマージする', () => {
    const result = buildHeatmapData(2026, [
      { date: '2026-03-15T00:00:00Z', total_min: 45 },
    ])

    const march15 = result.find((d) => d.date === '2026-03-15')
    expect(march15).toEqual({ date: '2026-03-15', count: 45, level: 2 })
  })
})

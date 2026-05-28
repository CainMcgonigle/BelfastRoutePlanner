import type { StopService } from './api'

interface CompressedBadge {
  label: string
  description: string
  operator: string
}

function consecutiveRanges(letters: string[]): string[] {
  if (letters.length === 0) return ['']
  const sorted = [...new Set(letters)].sort()
  const ranges: string[] = []
  let start = sorted[0]
  let end = sorted[0]
  for (let i = 1; i < sorted.length; i++) {
    if (sorted[i].charCodeAt(0) === end.charCodeAt(0) + 1) {
      end = sorted[i]
    } else {
      ranges.push(start === end ? start : `${start}-${end}`)
      start = sorted[i]
      end = sorted[i]
    }
  }
  ranges.push(start === end ? start : `${start}-${end}`)
  return ranges
}

export function compressServices(services: StopService[]): CompressedBadge[] {
  // Group by operator + numeric part of lineId
  const groups = new Map<string, { num: string; operator: string; suffixes: string[]; descriptions: string[] }>()

  for (const s of services) {
    const m = s.lineId.match(/^(N?)(\d+)([A-Z]*)$/)
    if (!m) {
      // doesn't match pattern (e.g. G1, G2) — show as-is
      const key = `${s.operator}:${s.lineId}`
      if (!groups.has(key)) groups.set(key, { num: s.lineId, operator: s.operator, suffixes: [], descriptions: [] })
      groups.get(key)!.descriptions.push(s.description)
      continue
    }
    const [, night, num, suffix] = m
    const key = `${s.operator}:${night}${num}`
    if (!groups.has(key)) groups.set(key, { num: `${night}${num}`, operator: s.operator, suffixes: [], descriptions: [] })
    groups.get(key)!.suffixes.push(suffix)
    groups.get(key)!.descriptions.push(s.description)
  }

  const badges: CompressedBadge[] = []
  for (const [, g] of groups) {
    const suffixes = g.suffixes.filter(s => s !== '')
    if (suffixes.length === 0) {
      badges.push({ label: g.num, description: g.descriptions[0] ?? '', operator: g.operator })
    } else {
      for (const range of consecutiveRanges(suffixes)) {
        badges.push({
          label: `${g.num}${range}`,
          description: g.descriptions.join(', '),
          operator: g.operator,
        })
      }
    }
  }

  return badges.sort((a, b) => a.label.localeCompare(b.label, undefined, { numeric: true }))
}

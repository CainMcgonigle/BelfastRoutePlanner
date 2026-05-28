'use client'

import { useState } from 'react'
import { planTrip } from '../lib/api'
import type { Itinerary, Leg } from '../lib/api'

function formatTime(epochMs: number): string {
  return new Date(epochMs).toLocaleTimeString('en-GB', {
    hour: '2-digit',
    minute: '2-digit',
  })
}

function formatDuration(seconds: number): string {
  const mins = Math.round(seconds / 60)
  if (mins < 60) return `${mins} min`
  const h = Math.floor(mins / 60)
  const m = mins % 60
  return m > 0 ? `${h}h ${m}m` : `${h}h`
}

const styles: Record<string, React.CSSProperties> = {
  container: {
    padding: '20px',
    backgroundColor: '#0f172a',
    minHeight: '100%',
    color: '#ffffff',
  },
  heading: {
    fontSize: '1.25rem',
    fontWeight: 700,
    marginBottom: '20px',
    color: '#f1f5f9',
    letterSpacing: '-0.01em',
  },
  label: {
    display: 'block',
    fontSize: '0.75rem',
    fontWeight: 600,
    color: '#94a3b8',
    marginBottom: '6px',
    textTransform: 'uppercase' as const,
    letterSpacing: '0.05em',
  },
  input: {
    width: '100%',
    padding: '10px 12px',
    backgroundColor: '#1e293b',
    border: '1px solid #334155',
    borderRadius: '8px',
    color: '#f1f5f9',
    fontSize: '0.875rem',
    marginBottom: '14px',
    outline: 'none',
  },
  button: {
    width: '100%',
    padding: '11px',
    backgroundColor: '#3b82f6',
    color: '#ffffff',
    border: 'none',
    borderRadius: '8px',
    fontSize: '0.875rem',
    fontWeight: 600,
    cursor: 'pointer',
    marginTop: '4px',
  },
  buttonDisabled: {
    opacity: 0.6,
    cursor: 'not-allowed',
  },
  error: {
    marginTop: '16px',
    padding: '10px 14px',
    backgroundColor: '#450a0a',
    border: '1px solid #991b1b',
    borderRadius: '8px',
    color: '#fca5a5',
    fontSize: '0.85rem',
  },
  itineraryList: {
    marginTop: '20px',
    display: 'flex',
    flexDirection: 'column' as const,
    gap: '12px',
  },
  itineraryCard: {
    backgroundColor: '#1e293b',
    border: '1px solid #334155',
    borderRadius: '10px',
    padding: '14px',
  },
  itineraryDuration: {
    fontSize: '0.8rem',
    color: '#94a3b8',
    marginBottom: '10px',
    fontWeight: 600,
  },
  legItem: {
    display: 'flex',
    alignItems: 'flex-start',
    gap: '10px',
    padding: '8px 0',
    borderTop: '1px solid #0f172a',
  },
  legMode: {
    flexShrink: 0,
    fontSize: '0.7rem',
    fontWeight: 700,
    padding: '2px 7px',
    borderRadius: '4px',
    backgroundColor: '#0f172a',
    color: '#60a5fa',
    textTransform: 'uppercase' as const,
    letterSpacing: '0.04em',
    marginTop: '1px',
  },
  legDetails: {
    flex: 1,
    fontSize: '0.82rem',
    color: '#cbd5e1',
    lineHeight: 1.4,
  },
  legRoute: {
    fontWeight: 600,
    color: '#f1f5f9',
  },
  legTime: {
    color: '#94a3b8',
    fontSize: '0.78rem',
  },
}

export default function TripPlanner() {
  const [from, setFrom] = useState('')
  const [to, setTo] = useState('')
  const [itineraries, setItineraries] = useState<Itinerary[]>([])
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  function parseCoord(value: string): [number, number] | null {
    const parts = value.split(',').map((s) => s.trim())
    if (parts.length !== 2) return null
    const lat = parseFloat(parts[0])
    const lon = parseFloat(parts[1])
    if (isNaN(lat) || isNaN(lon)) return null
    return [lat, lon]
  }

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    setError(null)
    setItineraries([])

    const fromCoord = parseCoord(from)
    const toCoord = parseCoord(to)

    if (!fromCoord) {
      setError('Invalid "From" coordinates. Use format: lat,lon')
      return
    }
    if (!toCoord) {
      setError('Invalid "To" coordinates. Use format: lat,lon')
      return
    }

    setLoading(true)
    try {
      const result = await planTrip(fromCoord[0], fromCoord[1], toCoord[0], toCoord[1])
      if (result.error) {
        setError(result.error.msg)
      } else {
        setItineraries(result.plan?.itineraries ?? [])
        if ((result.plan?.itineraries ?? []).length === 0) {
          setError('No routes found for these coordinates.')
        }
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to fetch route.')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div style={styles.container}>
      <h1 style={styles.heading}>Plan a Trip</h1>
      <form onSubmit={handleSubmit}>
        <label style={styles.label} htmlFor="from">
          From (lat, lon)
        </label>
        <input
          id="from"
          style={styles.input}
          type="text"
          placeholder="54.5973, -5.9301"
          value={from}
          onChange={(e) => setFrom(e.target.value)}
          required
        />

        <label style={styles.label} htmlFor="to">
          To (lat, lon)
        </label>
        <input
          id="to"
          style={styles.input}
          type="text"
          placeholder="54.6100, -5.9200"
          value={to}
          onChange={(e) => setTo(e.target.value)}
          required
        />

        <button
          type="submit"
          style={
            loading
              ? { ...styles.button, ...styles.buttonDisabled }
              : styles.button
          }
          disabled={loading}
        >
          {loading ? 'Searching...' : 'Find Routes'}
        </button>
      </form>

      {error && <div style={styles.error}>{error}</div>}

      {itineraries.length > 0 && (
        <div style={styles.itineraryList}>
          {itineraries.map((itin, i) => (
            <ItineraryCard key={i} itinerary={itin} index={i} />
          ))}
        </div>
      )}
    </div>
  )
}

function ItineraryCard({ itinerary, index }: { itinerary: Itinerary; index: number }) {
  return (
    <div style={styles.itineraryCard}>
      <div style={styles.itineraryDuration}>
        Option {index + 1} &mdash; {formatDuration(itinerary.duration)}
      </div>
      {itinerary.legs.map((leg, j) => (
        <LegRow key={j} leg={leg} isFirst={j === 0} />
      ))}
    </div>
  )
}

function LegRow({ leg, isFirst }: { leg: Leg; isFirst: boolean }) {
  return (
    <div style={isFirst ? { ...styles.legItem, borderTop: 'none' } : styles.legItem}>
      <span style={styles.legMode}>{leg.mode}</span>
      <div style={styles.legDetails}>
        {leg.routeShortName && (
          <span style={styles.legRoute}>Route {leg.routeShortName} &mdash; </span>
        )}
        <span>
          {leg.from.name} &rarr; {leg.to.name}
        </span>
        <br />
        <span style={styles.legTime}>
          {formatTime(leg.startTime)} &ndash; {formatTime(leg.endTime)}
        </span>
      </div>
    </div>
  )
}

const API_BASE = 'http://localhost:8081'

export interface Leg {
  mode: string
  startTime: number
  endTime: number
  from: { name: string; stopId: string; lat: number; lon: number }
  to: { name: string; stopId: string; lat: number; lon: number }
  legGeometry: { points: string }
  routeShortName: string
}

export interface Itinerary {
  duration: number
  legs: Leg[]
}

export interface OTPResponse {
  plan: { itineraries: Itinerary[] }
  error?: { msg: string }
}

export interface Stop {
  atco: string
  name: string
  lat: number
  lon: number
}

export async function planTrip(
  fromLat: number,
  fromLon: number,
  toLat: number,
  toLon: number,
  date?: string,
  time?: string
): Promise<OTPResponse> {
  const now = new Date()
  const resolvedDate =
    date ??
    now.toLocaleDateString('en-CA') // YYYY-MM-DD
  const resolvedTime =
    time ??
    now.toLocaleTimeString('en-GB', { hour: '2-digit', minute: '2-digit' }) // HH:MM

  const params = new URLSearchParams({
    fromLat: String(fromLat),
    fromLon: String(fromLon),
    toLat: String(toLat),
    toLon: String(toLon),
    date: resolvedDate,
    time: resolvedTime,
  })

  const res = await fetch(`${API_BASE}/api/plan?${params.toString()}`)
  if (!res.ok) {
    throw new Error(`Plan request failed: ${res.status} ${res.statusText}`)
  }
  return res.json() as Promise<OTPResponse>
}

export async function getStops(): Promise<Stop[]> {
  const res = await fetch(`${API_BASE}/api/stops`)
  if (!res.ok) {
    throw new Error(`Stops request failed: ${res.status} ${res.statusText}`)
  }
  return res.json() as Promise<Stop[]>
}

export async function getStop(stopId: string): Promise<Stop> {
  const res = await fetch(`${API_BASE}/api/stops/${encodeURIComponent(stopId)}`)
  if (!res.ok) {
    throw new Error(`Stop request failed: ${res.status} ${res.statusText}`)
  }
  return res.json() as Promise<Stop>
}

export async function getRoutesGeoJSON(): Promise<GeoJSON.FeatureCollection> {
  const res = await fetch(`${API_BASE}/api/routes/geojson`)
  if (!res.ok) {
    throw new Error(`Routes GeoJSON request failed: ${res.status} ${res.statusText}`)
  }
  return res.json() as Promise<GeoJSON.FeatureCollection>
}

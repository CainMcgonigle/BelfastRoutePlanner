'use client'

import { useEffect, useState, useCallback } from 'react'
import MapGL, { Source, Layer, Popup } from 'react-map-gl/maplibre'
import type { RasterLayer, CircleLayer, MapMouseEvent } from 'maplibre-gl'
import 'maplibre-gl/dist/maplibre-gl.css'
import { getStops, getRoutesGeoJSON, getStop, type Stop, type StopService } from '../lib/api'
import { compressServices } from '../lib/compressServices'

const BELFAST_CENTER = { latitude: 54.5973, longitude: -5.9301, zoom: 12 }

const osmLayer: RasterLayer = {
  id: 'osm-tiles',
  type: 'raster',
  source: 'osm',
  minzoom: 0,
  maxzoom: 19,
}

const CORRIDOR_COLOURS: Record<string, string> = {
  '1':  '#f97316', // orange
  '2':  '#dc2626', // red
  '3':  '#06b6d4', // turquoise
  '4':  '#7c2d3e', // brownish maroon
  '5':  '#0d9488', // teal
  '6':  '#4a7c59', // cactus green
  '7':  '#a3e635', // soft lime green
  '8':  '#7c3aed', // purple
  '9':  '#ec4899', // pink
  '10': '#eab308', // yellow
  '11': '#92400e', // brown
  '12': '#2d8fbd', // blue with green tint
  '13': '#f472b6', // light pink
  'G1': '#1e3a8a', // navy blue
  'G2': '#2563eb', // blue
}

const corridorMatchExpr = [
  'match', ['get', 'corridor'],
  ...Object.entries(CORRIDOR_COLOURS).flatMap(([k, v]) => [k, v]),
  '#ec4899', // default pink
]

const routeLineLayer = {
  id: 'routes',
  type: 'line' as const,
  source: 'routes',
  layout: {
    'line-cap': 'round' as const,
    'line-join': 'round' as const,
  },
  paint: {
    'line-color': corridorMatchExpr,
    'line-width': ['interpolate', ['linear'], ['zoom'], 10, 1.5, 15, 4],
    'line-opacity': 0.85,
    'line-offset': [
      '*',
      ['get', 'offsetIndex'],
      ['interpolate', ['linear'], ['zoom'], 10, 2.5, 15, 7],
    ],
  },
}

const stopCircleLayer: CircleLayer = {
  id: 'stops',
  type: 'circle',
  source: 'stops',
  paint: {
    'circle-radius': ['interpolate', ['linear'], ['zoom'], 10, 3, 15, 7],
    'circle-color': '#f97316',
    'circle-stroke-width': 1.5,
    'circle-stroke-color': '#fff',
  },
}

function stopsToGeoJSON(stops: Stop[]) {
  return {
    type: 'FeatureCollection' as const,
    features: stops
      .filter(s => s.lat !== 0 && s.lon !== 0)
      .map(s => ({
        type: 'Feature' as const,
        geometry: { type: 'Point' as const, coordinates: [s.lon, s.lat] },
        properties: { atco: s.atco, name: s.name || s.atco },
      })),
  }
}

type Service = StopService

interface PopupInfo {
  longitude: number
  latitude: number
  name: string
  atco: string
  services: Service[] | null
}

const OPERATOR_LABEL: Record<string, string> = {
  MET: 'Metro',
  GDR: 'Glider',
  ULB: 'Ulsterbus',
  GLD: 'Goldliner',
}

export default function Map() {
  const [stops, setStops] = useState<Stop[]>([])
  const [routes, setRoutes] = useState<GeoJSON.FeatureCollection>({ type: 'FeatureCollection', features: [] })
  const [popup, setPopup] = useState<PopupInfo | null>(null)

  useEffect(() => {
    getStops().then(setStops).catch(console.error)
    getRoutesGeoJSON().then(setRoutes).catch(console.error)
  }, [])

  const handleClick = useCallback((e: MapMouseEvent) => {
    const features = e.features
    if (!features || features.length === 0) {
      setPopup(null)
      return
    }
    const f = features[0]
    const [longitude, latitude] = (f.geometry as GeoJSON.Point).coordinates
    const atco: string = f.properties?.atco
    const name: string = f.properties?.name ?? atco

    setPopup({ longitude, latitude, name, atco, services: null })

    getStop(atco).then(detail => {
      setPopup(p => p?.atco === atco
        ? { ...p, services: (detail as unknown as { services: Service[] }).services ?? [] }
        : p
      )
    }).catch(console.error)
  }, [])

  return (
    <MapGL
      initialViewState={BELFAST_CENTER}
      style={{ width: '100%', height: '100%' }}
      mapStyle={{ version: 8, sources: {}, layers: [] }}
      interactiveLayerIds={['stops']}
      onClick={handleClick}
      cursor="auto"
    >
      <Source
        id="osm"
        type="raster"
        tiles={['https://tile.openstreetmap.org/{z}/{x}/{y}.png']}
        tileSize={256}
        attribution="&copy; OpenStreetMap contributors"
      >
        <Layer {...osmLayer} />
      </Source>

      <Source id="routes" type="geojson" data={routes}>
        <Layer {...routeLineLayer} />
      </Source>

      <Source id="stops" type="geojson" data={stopsToGeoJSON(stops)}>
        <Layer {...stopCircleLayer} />
      </Source>

      {popup && (
        <Popup
          longitude={popup.longitude}
          latitude={popup.latitude}
          anchor="bottom"
          onClose={() => setPopup(null)}
          closeButton
          maxWidth="260px"
        >
          <div style={{ color: '#0f172a', fontSize: 13, lineHeight: 1.6, minWidth: 180 }}>
            <div style={{ fontWeight: 600, marginBottom: 2 }}>{popup.name}</div>
            <div style={{ color: '#64748b', fontSize: 11, marginBottom: 8 }}>{popup.atco}</div>

            {popup.services === null && (
              <div style={{ color: '#94a3b8', fontSize: 11 }}>Loading services…</div>
            )}

            {popup.services && popup.services.length === 0 && (
              <div style={{ color: '#94a3b8', fontSize: 11 }}>No services found</div>
            )}

            {popup.services && popup.services.length > 0 && (
              <div>
                <div style={{ fontSize: 11, color: '#64748b', marginBottom: 4 }}>Services</div>
                <div style={{ display: 'flex', flexWrap: 'wrap', gap: 4 }}>
                  {compressServices(popup.services).map((b, i) => (
                    <span key={i} title={b.description} style={{
                      background: b.operator === 'GDR' ? '#d1fae5' : '#dbeafe',
                      color: b.operator === 'GDR' ? '#065f46' : '#1e40af',
                      borderRadius: 4,
                      padding: '1px 6px',
                      fontSize: 12,
                      fontWeight: 600,
                      cursor: 'default',
                    }}>
                      {b.label}
                    </span>
                  ))}
                </div>
              </div>
            )}
          </div>
        </Popup>
      )}
    </MapGL>
  )
}

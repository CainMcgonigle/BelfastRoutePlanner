'use client'

import { useEffect, useState, useCallback } from 'react'
import MapGL, { Source, Layer, Popup } from 'react-map-gl/maplibre'
import type { RasterLayer, CircleLayer, MapMouseEvent } from 'maplibre-gl'
import 'maplibre-gl/dist/maplibre-gl.css'
import { getStops, type Stop } from '../lib/api'

const BELFAST_CENTER = {
  latitude: 54.5973,
  longitude: -5.9301,
  zoom: 12,
}

const osmLayer: RasterLayer = {
  id: 'osm-tiles',
  type: 'raster',
  source: 'osm',
  minzoom: 0,
  maxzoom: 19,
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

interface PopupInfo {
  longitude: number
  latitude: number
  name: string
  atco: string
}

export default function Map() {
  const [stops, setStops] = useState<Stop[]>([])
  const [popup, setPopup] = useState<PopupInfo | null>(null)

  useEffect(() => {
    getStops().then(setStops).catch(console.error)
  }, [])

  const handleClick = useCallback((e: MapMouseEvent) => {
    const features = e.features
    if (!features || features.length === 0) {
      setPopup(null)
      return
    }
    const f = features[0]
    const [longitude, latitude] = (f.geometry as GeoJSON.Point).coordinates
    setPopup({
      longitude,
      latitude,
      name: f.properties?.name ?? f.properties?.atco,
      atco: f.properties?.atco,
    })
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
        >
          <div style={{ color: '#0f172a', fontSize: 13, lineHeight: 1.5 }}>
            <strong>{popup.name}</strong>
            <br />
            <span style={{ color: '#64748b', fontSize: 11 }}>{popup.atco}</span>
          </div>
        </Popup>
      )}
    </MapGL>
  )
}

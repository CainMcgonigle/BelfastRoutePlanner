'use client'

import { useRef } from 'react'
import MapGL, { Source, Layer } from 'react-map-gl/maplibre'
import type { RasterLayer } from 'maplibre-gl'
import 'maplibre-gl/dist/maplibre-gl.css'

const BELFAST_CENTER = {
  latitude: 54.5973,
  longitude: -5.9301,
  zoom: 12,
}

const osmRasterLayer: RasterLayer = {
  id: 'osm-tiles',
  type: 'raster',
  source: 'osm',
  minzoom: 0,
  maxzoom: 19,
}

export default function Map() {
  const mapRef = useRef(null)

  return (
    <MapGL
      ref={mapRef}
      initialViewState={BELFAST_CENTER}
      style={{ width: '100%', height: '100%' }}
      mapStyle={{
        version: 8,
        sources: {},
        layers: [],
      }}
    >
      <Source
        id="osm"
        type="raster"
        tiles={['https://tile.openstreetmap.org/{z}/{x}/{y}.png']}
        tileSize={256}
        attribution="&copy; <a href='https://www.openstreetmap.org/copyright'>OpenStreetMap</a> contributors"
      >
        <Layer {...osmRasterLayer} />
      </Source>
    </MapGL>
  )
}

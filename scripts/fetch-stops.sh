#!/usr/bin/env bash
set -euo pipefail

OUT="$(cd "$(dirname "$0")/.." && pwd)/data/stops.json"

echo "Querying OpenStreetMap Overpass API for NI bus stops..."

curl -s --progress-bar \
  --data-urlencode 'data=[out:json][timeout:60];
area["name"="Northern Ireland"]["admin_level"="4"]->.ni;
(
  node["highway"="bus_stop"](area.ni);
  node["public_transport"="stop_position"]["bus"="yes"](area.ni);
);
out body;' \
  "https://overpass-api.de/api/interpreter" \
  -o "$OUT"

COUNT=$(python3 -c "import json,sys; d=json.load(open('$OUT')); print(len(d['elements']))" 2>/dev/null || echo "unknown")
echo "Saved $COUNT stops to $OUT"

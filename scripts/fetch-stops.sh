#!/usr/bin/env bash
set -euo pipefail

OUT="$(cd "$(dirname "$0")/.." && pwd)/data/stops.json"

QUERY='[out:json][timeout:60];
(
  node["highway"="bus_stop"](54.0,-8.2,55.4,-5.4);
  node["public_transport"="stop_position"]["bus"="yes"](54.0,-8.2,55.4,-5.4);
);
out body;'

MIRRORS=(
  "https://overpass-api.de/api/interpreter"
  "https://overpass.kumi.systems/api/interpreter"
  "https://overpass.openstreetmap.ru/api/interpreter"
)

echo "Querying OpenStreetMap Overpass API for NI bus stops..."

SUCCESS=false
for MIRROR in "${MIRRORS[@]}"; do
  echo "Trying $MIRROR ..."
  HTTP_CODE=$(curl -s -o "$OUT" -w "%{http_code}" \
    -A "BelfastRoutePlanner/1.0 (github.com/CainMcgonigle/BelfastRoutePlanner)" \
    -X POST "$MIRROR" \
    --data-urlencode "data=$QUERY")

  if [ "$HTTP_CODE" = "200" ] && python3 -c "import json; json.load(open('$OUT'))" 2>/dev/null; then
    SUCCESS=true
    break
  else
    echo "  Failed (HTTP $HTTP_CODE), trying next mirror..."
  fi
done

if [ "$SUCCESS" = false ]; then
  echo ""
  echo "All mirrors failed. Try manually in a browser:"
  echo "  https://overpass-api.de/api/interpreter"
  echo "Paste the query, save the JSON result to: $OUT"
  exit 1
fi

COUNT=$(python3 -c "import json; d=json.load(open('$OUT')); print(len(d.get('elements', [])))")
echo "Saved $COUNT stops to $OUT"
echo ""

# Show sample tags to confirm ATCO code field name
echo "Sample tags from first stop:"
python3 -c "
import json
els = json.load(open('$OUT')).get('elements', [])
if els:
    print(json.dumps(els[0].get('tags', {}), indent=2))
else:
    print('No elements found - check bounding box')
"

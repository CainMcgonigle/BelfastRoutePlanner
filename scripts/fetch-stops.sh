#!/usr/bin/env bash
set -euo pipefail

OUT="$(cd "$(dirname "$0")/.." && pwd)/data/stops.json"
QUERY_FILE="$(mktemp /tmp/overpass-query-XXXXXX.txt)"

# Northern Ireland bounding box: south,west,north,east
cat > "$QUERY_FILE" <<'OVERPASS'
[out:json][timeout:60];
(
  node["highway"="bus_stop"](54.0,-8.2,55.4,-5.4);
  node["public_transport"="stop_position"]["bus"="yes"](54.0,-8.2,55.4,-5.4);
);
out body;
OVERPASS

echo "Querying OpenStreetMap Overpass API for NI bus stops..."
curl -s --progress-bar \
  -X POST \
  -H "Content-Type: application/x-www-form-urlencoded" \
  --data-urlencode "data@$QUERY_FILE" \
  "https://overpass-api.de/api/interpreter" \
  -o "$OUT"

rm -f "$QUERY_FILE"

# count results
if command -v python3 &>/dev/null; then
  COUNT=$(python3 -c "import json; d=json.load(open('$OUT')); print(len(d.get('elements', [])))")
  echo "Saved $COUNT stops to $OUT"
else
  echo "Saved to $OUT"
fi

# show a sample of tags from first stop to help debug ATCO matching
echo ""
echo "Sample tags from first stop:"
python3 -c "
import json
d = json.load(open('$OUT'))
els = d.get('elements', [])
if els:
    print(json.dumps(els[0].get('tags', {}), indent=2))
else:
    print('No elements found')
" 2>/dev/null || true

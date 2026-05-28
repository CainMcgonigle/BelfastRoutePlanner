#!/usr/bin/env python3
import json, urllib.request

routes = json.loads(urllib.request.urlopen("http://localhost:8081/api/routes").read())
stops = json.loads(urllib.request.urlopen("http://localhost:8081/api/stops").read())

stop_map = {s['atco']: s for s in stops}

for r in routes:
    line = r['lineId'].strip()
    if not line.startswith('9'):
        continue
    journeys = r.get('journeys', [])
    if not journeys:
        continue
    stop_times = journeys[0].get('stopTimes', [])
    total = len(stop_times)
    with_coords = sum(1 for st in stop_times if stop_map.get(st['atco'], {}).get('lat', 0) != 0)
    print(f"{line:4s}  {with_coords}/{total} stops with coords  {r['description']}")

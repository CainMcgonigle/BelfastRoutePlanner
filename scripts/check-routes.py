#!/usr/bin/env python3
import json, urllib.request
from collections import defaultdict

routes = json.loads(urllib.request.urlopen("http://localhost:8081/api/routes").read())

by_operator = defaultdict(list)
for r in routes:
    by_operator[r['operator']].append(r)

for op, op_routes in sorted(by_operator.items()):
    lines = sorted(set(r['lineId'].strip() for r in op_routes))
    print(f"{op}: {len(op_routes)} routes — lines: {', '.join(lines[:20])}{'...' if len(lines) > 20 else ''}")

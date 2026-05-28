#!/usr/bin/env python3
import json, sys

path = sys.argv[1] if len(sys.argv) > 1 else "../data/stops.json"
elements = json.load(open(path)).get("elements", [])

print("Sample naptan:AtcoCode values:")
count = 0
for el in elements:
    tags = el.get("tags", {})
    if "naptan:AtcoCode" in tags:
        print(f"  naptan:AtcoCode={tags['naptan:AtcoCode']}  ref={tags.get('ref', '-')}  name={tags.get('name', '-')}")
        count += 1
        if count >= 10:
            break

print("\nSample ref values (stops without naptan:AtcoCode):")
count = 0
for el in elements:
    tags = el.get("tags", {})
    if "ref" in tags and "naptan:AtcoCode" not in tags:
        print(f"  ref={tags['ref']}  name={tags.get('name', '-')}")
        count += 1
        if count >= 10:
            break

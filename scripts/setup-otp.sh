#!/usr/bin/env bash
set -euo pipefail

GRAPH_DIR="$(cd "$(dirname "$0")/.." && pwd)/otp/graphs/ni"
OTP_DIR="$(cd "$(dirname "$0")/.." && pwd)/otp"

OSM_URL="https://download.geofabrik.de/europe/great-britain/northern-ireland-latest.osm.pbf"
OSM_FILE="$GRAPH_DIR/northern-ireland-latest.osm.pbf"
GTFS_FILE="$GRAPH_DIR/translink-ni.gtfs.zip"
OTP_JAR="$OTP_DIR/otp.jar"

mkdir -p "$GRAPH_DIR"

# --- OTP jar ---
if [ ! -f "$OTP_JAR" ]; then
    echo "Fetching latest OTP release URL..."
    DOWNLOAD_URL=$(curl -s https://api.github.com/repos/opentripplanner/OpenTripPlanner/releases/latest \
        | grep "browser_download_url" \
        | grep "shaded.jar" \
        | cut -d '"' -f 4)

    if [ -z "$DOWNLOAD_URL" ]; then
        echo "ERROR: could not resolve OTP download URL. Check your internet connection."
        exit 1
    fi

    echo "Downloading OTP from $DOWNLOAD_URL ..."
    curl -L -o "$OTP_JAR" "$DOWNLOAD_URL"
    echo "OTP jar saved to $OTP_JAR"
else
    echo "OTP jar already present, skipping."
fi

# --- OSM data ---
if [ ! -f "$OSM_FILE" ]; then
    echo "Downloading Northern Ireland OSM data..."
    curl -L -o "$OSM_FILE" "$OSM_URL"
    echo "OSM data saved."
else
    echo "OSM data already present, skipping."
fi

# --- GTFS data ---
if [ ! -f "$GTFS_FILE" ]; then
    echo ""
    echo "GTFS data not found at $GTFS_FILE"
    echo ""
    echo "Download the Translink NI GTFS zip from OpenDataNI:"
    echo "  https://www.opendatani.gov.uk  (search: Translink GTFS)"
    echo ""
    echo "Then either:"
    echo "  1. Paste the direct download URL when prompted below"
    echo "  2. Or manually place the file at: $GTFS_FILE"
    echo ""
    read -rp "Paste GTFS download URL (or press Enter to skip and place file manually): " GTFS_URL

    if [ -n "$GTFS_URL" ]; then
        curl -L -o "$GTFS_FILE" "$GTFS_URL"
        echo "GTFS data saved."
    else
        echo "Skipping GTFS download. Place the file manually before building the graph."
        echo "Expected path: $GTFS_FILE"
        exit 0
    fi
else
    echo "GTFS data already present, skipping."
fi

# --- Build OTP graph ---
echo ""
echo "All data present. Building OTP graph (this takes a few minutes)..."
java -Xmx2G -jar "$OTP_JAR" --build --save "$GRAPH_DIR"

echo ""
echo "Graph built successfully."
echo "Start OTP with:"
echo "  java -Xmx2G -jar $OTP_JAR --load --serve $GRAPH_DIR"
echo ""
echo "Or via Docker Compose from the repo root:"
echo "  docker compose up"

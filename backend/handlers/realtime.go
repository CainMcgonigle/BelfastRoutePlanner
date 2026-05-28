package handlers

import (
	"encoding/json"
	"net/http"
)

// placeholder until OpenDataNI publishes a GTFS-RT feed
func GetRealtimeArrivals(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "realtime not yet available"})
}

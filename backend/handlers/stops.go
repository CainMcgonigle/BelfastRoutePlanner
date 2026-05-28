package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// placeholder until GTFS stops are loaded
func GetStops(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode([]interface{}{})
}

func GetStop(w http.ResponseWriter, r *http.Request) {
	stopID := chi.URLParam(r, "stopID")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"stopId": stopID})
}

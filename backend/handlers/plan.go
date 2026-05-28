package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"transit-ni/services"
)

func PlanTrip(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	fromLat, _ := strconv.ParseFloat(q.Get("fromLat"), 64)
	fromLon, _ := strconv.ParseFloat(q.Get("fromLon"), 64)
	toLat, _ := strconv.ParseFloat(q.Get("toLat"), 64)
	toLon, _ := strconv.ParseFloat(q.Get("toLon"), 64)

	date := q.Get("date")
	if date == "" {
		date = time.Now().Format("01-02-2006")
	}
	tripTime := q.Get("time")
	if tripTime == "" {
		tripTime = time.Now().Format("3:04pm")
	}

	result, err := services.PlanTrip(services.PlanRequest{
		FromLat: fromLat, FromLon: fromLon,
		ToLat: toLat, ToLon: toLon,
		Date: date, Time: tripTime,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

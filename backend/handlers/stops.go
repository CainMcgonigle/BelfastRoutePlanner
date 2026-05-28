package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"transit-ni/services"
)

func GetStops(w http.ResponseWriter, r *http.Request) {
	data := services.GetCIFData()
	stops := make([]*services.CIFStop, 0, len(data.Stops))
	for _, s := range data.Stops {
		stops = append(stops, s)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stops)
}

type stopDetail struct {
	*services.CIFStop
	Services []services.StopService `json:"services"`
}

func GetStop(w http.ResponseWriter, r *http.Request) {
	stopID := chi.URLParam(r, "stopID")
	data := services.GetCIFData()
	stop, ok := data.Stops[stopID]
	if !ok {
		http.Error(w, "stop not found", http.StatusNotFound)
		return
	}
	services_ := data.StopRoutes[stopID]
	if services_ == nil {
		services_ = []services.StopService{}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stopDetail{stop, services_})
}

func GetRoutes(w http.ResponseWriter, r *http.Request) {
	data := services.GetCIFData()
	routes := make([]*services.CIFRoute, 0, len(data.Routes))
	for _, r := range data.Routes {
		routes = append(routes, r)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(routes)
}

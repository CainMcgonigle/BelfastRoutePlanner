package handlers

import (
	"encoding/json"
	"net/http"

	"transit-ni/services"
)

type geojsonFeature struct {
	Type       string                 `json:"type"`
	Properties map[string]interface{} `json:"properties"`
	Geometry   geojsonGeometry        `json:"geometry"`
}

type geojsonGeometry struct {
	Type        string      `json:"type"`
	Coordinates interface{} `json:"coordinates"`
}

type geojsonCollection struct {
	Type     string           `json:"type"`
	Features []geojsonFeature `json:"features"`
}

func GetRoutesGeoJSON(w http.ResponseWriter, r *http.Request) {
	data := services.GetCIFData()
	features := []geojsonFeature{}

	for _, route := range data.Routes {
		// use one representative journey per route to draw the line
		if len(route.Journeys) == 0 {
			continue
		}
		coords := journeyCoords(route.Journeys[0], data.Stops)
		if len(coords) < 2 {
			continue
		}
		features = append(features, geojsonFeature{
			Type: "Feature",
			Properties: map[string]interface{}{
				"lineId":      route.LineID,
				"description": route.Description,
				"operator":    route.Operator,
				"direction":   route.Direction,
			},
			Geometry: geojsonGeometry{
				Type:        "LineString",
				Coordinates: coords,
			},
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(geojsonCollection{Type: "FeatureCollection", Features: features})
}

// journeyCoords returns [lon, lat] pairs for stops in journey order that have coordinates.
// Consecutive stops without coordinates break the line into separate segments.
func journeyCoords(journey services.CIFJourney, stops map[string]*services.CIFStop) [][2]float64 {
	coords := [][2]float64{}
	for _, st := range journey.StopTimes {
		stop, ok := stops[st.ATCO]
		if !ok || stop.Lat == 0 {
			continue
		}
		coords = append(coords, [2]float64{stop.Lon, stop.Lat})
	}
	return coords
}

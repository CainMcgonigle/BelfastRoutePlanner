package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"unicode"

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
	snapRoads := services.OSRMAvailable()
	features := []geojsonFeature{}

	for _, route := range data.Routes {
		if len(route.Journeys) == 0 {
			continue
		}
		coords := journeyCoords(route.Journeys[0], data.Stops)
		if len(coords) < 2 {
			continue
		}
		if snapRoads {
			coords = services.SnapToRoads(coords)
		}
		features = append(features, geojsonFeature{
			Type: "Feature",
			Properties: map[string]interface{}{
				"lineId":      route.LineID,
				"description": route.Description,
				"operator":    route.Operator,
				"direction":   route.Direction,
				"corridor":    corridor(route.LineID),
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

// corridor extracts the route prefix used for colour coding:
// "9A" → "9", "12B" → "12", "G1" → "G1", "N9" → "N9"
func corridor(lineID string) string {
	s := strings.TrimSpace(lineID)
	if s == "" {
		return "other"
	}
	// Glider and Night services keep their full prefix
	if strings.HasPrefix(s, "G") || strings.HasPrefix(s, "N") {
		return s[:2]
	}
	// Extract leading digits
	end := 0
	for _, c := range s {
		if !unicode.IsDigit(c) {
			break
		}
		end++
	}
	if end == 0 {
		return s
	}
	return s[:end]
}

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

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
	routeFeatures := services.GetRouteFeatures()
	features := make([]geojsonFeature, 0, len(routeFeatures))

	for _, rf := range routeFeatures {
		features = append(features, geojsonFeature{
			Type: "Feature",
			Properties: map[string]interface{}{
				"lineId":      rf.Props.LineID,
				"description": rf.Props.Description,
				"operator":    rf.Props.Operator,
				"direction":   rf.Props.Direction,
				"corridor":    rf.Props.Corridor,
				"offsetIndex": rf.Props.OffsetIndex,
			},
			Geometry: geojsonGeometry{
				Type:        "LineString",
				Coordinates: rf.Coords,
			},
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(geojsonCollection{Type: "FeatureCollection", Features: features})
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

func corridor(lineID string) string {
	s := strings.TrimSpace(lineID)
	if s == "" {
		return "other"
	}
	if strings.HasPrefix(s, "G") || strings.HasPrefix(s, "N") {
		return s[:2]
	}
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

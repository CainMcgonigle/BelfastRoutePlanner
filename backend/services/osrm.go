package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

var osrmBase string

func init() {
	osrmBase = os.Getenv("OSRM_BASE")
	if osrmBase == "" {
		osrmBase = "http://localhost:5000"
	}
}

type osrmRoute struct {
	Geometry struct {
		Coordinates [][2]float64 `json:"coordinates"`
	} `json:"geometry"`
}

type osrmResponse struct {
	Code   string      `json:"code"`
	Routes []osrmRoute `json:"routes"`
}

// SnapToRoads takes an ordered list of [lon,lat] waypoints and returns
// road-snapped coordinates. Falls back to the original points on error.
func SnapToRoads(waypoints [][2]float64) [][2]float64 {
	if len(waypoints) < 2 {
		return waypoints
	}

	coords := make([]string, len(waypoints))
	for i, wp := range waypoints {
		coords[i] = fmt.Sprintf("%f,%f", wp[0], wp[1])
	}

	url := fmt.Sprintf("%s/route/v1/driving/%s?overview=full&geometries=geojson",
		osrmBase, strings.Join(coords, ";"))

	resp, err := http.Get(url)
	if err != nil {
		return waypoints
	}
	defer resp.Body.Close()

	var result osrmResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return waypoints
	}
	if result.Code != "Ok" || len(result.Routes) == 0 {
		return waypoints
	}

	return result.Routes[0].Geometry.Coordinates
}

// OSRMAvailable returns true if the OSRM service is reachable.
func OSRMAvailable() bool {
	resp, err := http.Get(osrmBase + "/health")
	if err != nil {
		return false
	}
	resp.Body.Close()
	return resp.StatusCode == 200
}

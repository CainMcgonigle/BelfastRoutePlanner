package services

import (
	"encoding/json"
	"os"
)

type translinkFeature struct {
	Properties struct {
		AtcoCode   string `json:"AtcoCode"`
		CommonName string `json:"CommonName"`
		Status     string `json:"Status"`
	} `json:"properties"`
	Geometry struct {
		Coordinates [2]float64 `json:"coordinates"` // [lon, lat]
	} `json:"geometry"`
}

type translinkGeoJSON struct {
	Features []translinkFeature `json:"features"`
}

// EnrichStopsFromGeoJSON loads the Translink bus stop GeoJSON from OpenDataNI
// and fills in Name/Lat/Lon for any CIFStop whose ATCO code matches.
func EnrichStopsFromGeoJSON(data *CIFData, path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	var geojson translinkGeoJSON
	if err := json.NewDecoder(f).Decode(&geojson); err != nil {
		return err
	}

	for _, feature := range geojson.Features {
		atco := feature.Properties.AtcoCode
		if atco == "" || feature.Properties.Status != "active" {
			continue
		}
		stop, ok := data.Stops[atco]
		if !ok {
			continue
		}
		stop.Lat = feature.Geometry.Coordinates[1]
		stop.Lon = feature.Geometry.Coordinates[0]
		stop.Name = feature.Properties.CommonName
	}

	return nil
}

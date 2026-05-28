package services

import (
	"encoding/json"
	"os"
)

type osmNode struct {
	ID   int64              `json:"id"`
	Lat  float64            `json:"lat"`
	Lon  float64            `json:"lon"`
	Tags map[string]string  `json:"tags"`
}

type osmResponse struct {
	Elements []osmNode `json:"elements"`
}

// EnrichStopsFromOSM loads the Overpass JSON file and fills in Name/Lat/Lon
// for any CIFStop whose ATCO code matches a "naptan:AtcoCode" or "ref" tag.
func EnrichStopsFromOSM(data *CIFData, path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	var resp osmResponse
	if err := json.NewDecoder(f).Decode(&resp); err != nil {
		return err
	}

	matched := 0
	for _, node := range resp.Elements {
		atco := node.Tags["naptan:AtcoCode"]
		if atco == "" {
			atco = node.Tags["ref"]
		}
		if atco == "" {
			continue
		}
		stop, ok := data.Stops[atco]
		if !ok {
			continue
		}
		stop.Lat = node.Lat
		stop.Lon = node.Lon
		if name := node.Tags["name"]; name != "" {
			stop.Name = name
		} else if name := node.Tags["naptan:CommonName"]; name != "" {
			stop.Name = name
		}
		matched++
	}

	return nil
}

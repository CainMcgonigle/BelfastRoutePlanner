package services

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
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
		atco := resolveATCO(node.Tags)
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

// resolveATCO tries naptan:AtcoCode first, then constructs from ref by
// parsing as an integer and formatting as 700 + 9-digit zero-padded number.
// Refs like "95" or "900;901" (route refs, not stop refs) are skipped.
func resolveATCO(tags map[string]string) string {
	if v := tags["naptan:AtcoCode"]; v != "" {
		return v
	}
	ref := tags["ref"]
	if ref == "" || strings.ContainsAny(ref, ";, ") {
		return ""
	}
	n, err := strconv.ParseInt(ref, 10, 64)
	if err != nil || n < 100 || n > 999999999 {
		return ""
	}
	return fmt.Sprintf("700%09d", n)
}

package services

import (
	"math"
	"sort"
	"strings"
	"sync"
	"time"
)

type RouteFeatureProps struct {
	LineID      string `json:"lineId"`
	Corridor    string `json:"corridor"`
	Operator    string `json:"operator"`
	Description string `json:"description"`
	Direction   string `json:"direction"`
	OffsetIndex int    `json:"offsetIndex"`
}

type RouteFeature struct {
	Coords [][2]float64
	Props  RouteFeatureProps
}

var (
	routeMu       sync.RWMutex
	cachedFeatures []RouteFeature
	routesReady   bool
)

// WarmRouteCache pre-computes route features in the background,
// waiting for OSRM to become available first.
func WarmRouteCache() {
	go func() {
		for i := 0; i < 30; i++ {
			if OSRMAvailable() {
				break
			}
			time.Sleep(2 * time.Second)
		}
		features := computeRouteFeatures(GetCIFData())
		routeMu.Lock()
		cachedFeatures = features
		routesReady = true
		routeMu.Unlock()
	}()
}

// GetRouteFeatures returns the cached route features, computing them
// synchronously on first call if the background warm hasn't finished yet.
func GetRouteFeatures() []RouteFeature {
	routeMu.RLock()
	if routesReady {
		f := cachedFeatures
		routeMu.RUnlock()
		return f
	}
	routeMu.RUnlock()

	// Not ready yet — compute synchronously (OSRM may not be available)
	features := computeRouteFeatures(GetCIFData())
	routeMu.Lock()
	if !routesReady {
		cachedFeatures = features
		routesReady = true
	}
	routeMu.Unlock()
	return cachedFeatures
}

type routeWithGeom struct {
	key       string
	corridor  string
	lineID    string
	operator  string
	desc      string
	direction string
	coords    [][2]float64
}

func computeRouteFeatures(data *CIFData) []RouteFeature {
	// Collect one representative journey per route with (optionally) snapped coords
	snap := OSRMAvailable()
	var routes []routeWithGeom
	for key, route := range data.Routes {
		if len(route.Journeys) == 0 {
			continue
		}
		coords := journeyStopCoords(route.Journeys[0], data.Stops)
		if len(coords) < 2 {
			continue
		}
		if snap {
			coords = SnapToRoads(coords)
		}
		routes = append(routes, routeWithGeom{
			key:       key,
			corridor:  CorridorFromLine(route.LineID),
			lineID:    strings.TrimSpace(route.LineID),
			operator:  route.Operator,
			desc:      route.Description,
			direction: route.Direction,
			coords:    coords,
		})
	}

	// Sort routes for consistent offset ordering (by corridor then lineID)
	sort.Slice(routes, func(i, j int) bool {
		if routes[i].corridor != routes[j].corridor {
			return routes[i].corridor < routes[j].corridor
		}
		return routes[i].key < routes[j].key
	})

	// Build edge → route-key list
	type edgeKey [4]int64
	round := func(v float64) int64 { return int64(math.Round(v * 1e5)) }
	makeEdge := func(a, b [2]float64) edgeKey {
		a0, a1 := round(a[0]), round(a[1])
		b0, b1 := round(b[0]), round(b[1])
		if a0 < b0 || (a0 == b0 && a1 < b1) {
			return edgeKey{a0, a1, b0, b1}
		}
		return edgeKey{b0, b1, a0, a1}
	}

	edgeRoutes := make(map[edgeKey][]string)
	edgeSeen := make(map[edgeKey]map[string]bool)
	for _, r := range routes {
		for i := 0; i < len(r.coords)-1; i++ {
			e := makeEdge(r.coords[i], r.coords[i+1])
			if edgeSeen[e] == nil {
				edgeSeen[e] = make(map[string]bool)
			}
			if !edgeSeen[e][r.key] {
				edgeSeen[e][r.key] = true
				edgeRoutes[e] = append(edgeRoutes[e], r.key)
			}
		}
	}

	// Split each route into sub-segments wherever offsetIndex changes
	keyIndex := make(map[string]int, len(routes))
	for i, r := range routes {
		keyIndex[r.key] = i
	}

	var features []RouteFeature

	for _, r := range routes {
		if len(r.coords) < 2 {
			continue
		}

		offsetAt := func(i int) int {
			e := makeEdge(r.coords[i], r.coords[i+1])
			list := edgeRoutes[e]
			pos := 0
			for j, k := range list {
				if k == r.key {
					pos = j
					break
				}
			}
			return pos - (len(list)-1)/2
		}

		segStart := 0
		curOffset := offsetAt(0)

		flush := func(end int, offset int) {
			seg := r.coords[segStart : end+1]
			if len(seg) >= 2 {
				features = append(features, RouteFeature{
					Coords: seg,
					Props: RouteFeatureProps{
						LineID:      r.lineID,
						Corridor:    r.corridor,
						Operator:    r.operator,
						Description: r.desc,
						Direction:   r.direction,
						OffsetIndex: offset,
					},
				})
			}
		}

		for i := 1; i < len(r.coords)-1; i++ {
			newOffset := offsetAt(i)
			if newOffset != curOffset {
				flush(i, curOffset)
				segStart = i
				curOffset = newOffset
			}
		}
		flush(len(r.coords)-1, curOffset)
	}

	return features
}

func journeyStopCoords(journey CIFJourney, stops map[string]*CIFStop) [][2]float64 {
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

// CorridorFromLine extracts the route prefix used for colour coding.
func CorridorFromLine(lineID string) string {
	s := strings.TrimSpace(lineID)
	if s == "" {
		return "other"
	}
	if strings.HasPrefix(s, "G") || strings.HasPrefix(s, "N") {
		if len(s) >= 2 {
			return s[:2]
		}
		return s
	}
	end := 0
	for _, c := range s {
		if c < '0' || c > '9' {
			break
		}
		end++
	}
	if end == 0 {
		return s
	}
	return s[:end]
}

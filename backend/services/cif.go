package services

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type CIFStop struct {
	ATCO string  `json:"atco"`
	Name string  `json:"name"`
	Lat  float64 `json:"lat"`
	Lon  float64 `json:"lon"`
}

type CIFStopTime struct {
	ATCO      string `json:"atco"`
	Arrival   string `json:"arrival"`
	Departure string `json:"departure"`
}

type CIFJourney struct {
	Operator   string        `json:"operator"`
	Code       string        `json:"code"`
	RouteID    string        `json:"routeId"`
	Direction  string        `json:"direction"`
	StartDate  string        `json:"startDate"`
	EndDate    string        `json:"endDate"`
	DaysOfWeek string        `json:"daysOfWeek"`
	StopTimes  []CIFStopTime `json:"stopTimes"`
}

type CIFRoute struct {
	Operator    string       `json:"operator"`
	LineID      string       `json:"lineId"`
	Direction   string       `json:"direction"`
	Description string       `json:"description"`
	Journeys    []CIFJourney `json:"journeys"`
}

type CIFData struct {
	Routes map[string]*CIFRoute // keyed by "operator:lineID:direction"
	Stops  map[string]*CIFStop  // keyed by ATCO code
}

var (
	cifOnce sync.Once
	cifData *CIFData
)

func LoadCIFDir(dir string) (*CIFData, error) {
	data := &CIFData{
		Routes: make(map[string]*CIFRoute),
		Stops:  make(map[string]*CIFStop),
	}

	files, err := filepath.Glob(filepath.Join(dir, "*.cif"))
	if err != nil {
		return nil, err
	}

	for _, f := range files {
		if err := parseCIF(f, data); err != nil {
			return nil, err
		}
	}

	return data, nil
}

func parseCIF(path string, data *CIFData) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var currentRoute *CIFRoute
	var currentJourney *CIFJourney

	for scanner.Scan() {
		line := scanner.Text()
		if len(line) < 2 {
			continue
		}

		switch line[:2] {
		case "QD":
			// Route description: QD + N/D + operator(4) + lineID(4) + direction(1) + description
			if len(line) < 12 {
				continue
			}
			operator := strings.TrimSpace(line[3:7])
			lineID := strings.TrimSpace(line[7:11])
			direction := string(line[11])
			description := ""
			if len(line) > 12 {
				description = strings.TrimSpace(line[12:])
			}
			key := operator + ":" + lineID + ":" + direction
			if _, exists := data.Routes[key]; !exists {
				data.Routes[key] = &CIFRoute{
					Operator:    operator,
					LineID:      lineID,
					Direction:   direction,
					Description: description,
				}
			}
			currentRoute = data.Routes[key]
			currentJourney = nil

		case "QS":
			// Journey header: QS + N/D + operator(4) + code(6) + startDate(8) + endDate(8) + daysOfWeek(7) + ...
			if len(line) < 45 || currentRoute == nil {
				continue
			}
			operator := strings.TrimSpace(line[3:7])
			code := strings.TrimSpace(line[7:13])
			startDate := line[13:21]
			endDate := line[21:29]
			daysOfWeek := line[29:36]
			routeID := ""
			if len(line) > 42 {
				routeID = strings.TrimSpace(line[38:42])
			}
			currentJourney = &CIFJourney{
				Operator:   operator,
				Code:       code,
				RouteID:    routeID,
				Direction:  currentRoute.Direction,
				StartDate:  startDate,
				EndDate:    endDate,
				DaysOfWeek: daysOfWeek,
			}
			currentRoute.Journeys = append(currentRoute.Journeys, *currentJourney)
			// point to the appended copy
			currentJourney = &currentRoute.Journeys[len(currentRoute.Journeys)-1]

		case "QO":
			// Origin: QO + atco(12) + departure(4)
			if len(line) < 18 || currentJourney == nil {
				continue
			}
			atco := line[2:14]
			depart := line[14:18]
			ensureStop(data, atco)
			currentJourney.StopTimes = append(currentJourney.StopTimes, CIFStopTime{
				ATCO:      atco,
				Departure: depart,
				Arrival:   depart,
			})

		case "QI":
			// Intermediate: QI + atco(12) + arrival(4) + departure(4)
			if len(line) < 22 || currentJourney == nil {
				continue
			}
			atco := line[2:14]
			arrival := line[14:18]
			departure := line[18:22]
			ensureStop(data, atco)
			currentJourney.StopTimes = append(currentJourney.StopTimes, CIFStopTime{
				ATCO:      atco,
				Arrival:   arrival,
				Departure: departure,
			})

		case "QT":
			// Terminus: QT + atco(12) + arrival(4)
			if len(line) < 18 || currentJourney == nil {
				continue
			}
			atco := line[2:14]
			arrival := line[14:18]
			ensureStop(data, atco)
			currentJourney.StopTimes = append(currentJourney.StopTimes, CIFStopTime{
				ATCO:    atco,
				Arrival: arrival,
			})
			currentJourney = nil
		}
	}

	return scanner.Err()
}

func ensureStop(data *CIFData, atco string) {
	if _, exists := data.Stops[atco]; !exists {
		data.Stops[atco] = &CIFStop{ATCO: atco}
	}
}

func GetCIFData() *CIFData {
	cifOnce.Do(func() {
		dir := os.Getenv("CIF_DIR")
		if dir == "" {
			dir = "../data/cif"
		}
		var err error
		cifData, err = LoadCIFDir(dir)
		if err != nil {
			cifData = &CIFData{
				Routes: make(map[string]*CIFRoute),
				Stops:  make(map[string]*CIFStop),
			}
			return
		}
		osmPath := os.Getenv("OSM_STOPS_FILE")
		if osmPath == "" {
			osmPath = "../data/stops.json"
		}
		if _, err := os.Stat(osmPath); err == nil {
			_ = EnrichStopsFromOSM(cifData, osmPath)
		}
	})
	return cifData
}

package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const otpBase = "http://localhost:8080/otp/routers/default"

type PlanRequest struct {
	FromLat, FromLon float64
	ToLat, ToLon     float64
	Date, Time       string
}

type OTPResponse struct {
	Plan struct {
		Itineraries []Itinerary `json:"itineraries"`
	} `json:"plan"`
	Error *struct {
		Message string `json:"msg"`
	} `json:"error"`
}

type Itinerary struct {
	Duration int   `json:"duration"`
	Legs     []Leg `json:"legs"`
}

type Leg struct {
	Mode           string   `json:"mode"`
	StartTime      int64    `json:"startTime"`
	EndTime        int64    `json:"endTime"`
	From           Place    `json:"from"`
	To             Place    `json:"to"`
	LegGeometry    Geometry `json:"legGeometry"`
	RouteShortName string   `json:"routeShortName"`
}

type Place struct {
	Name      string  `json:"name"`
	StopID    string  `json:"stopId"`
	Lat       float64 `json:"lat"`
	Lon       float64 `json:"lon"`
	Departure int64   `json:"departure"`
	Arrival   int64   `json:"arrival"`
}

type Geometry struct {
	Points string `json:"points"`
}

func PlanTrip(req PlanRequest) (*OTPResponse, error) {
	params := url.Values{}
	params.Set("fromPlace", fmt.Sprintf("%f,%f", req.FromLat, req.FromLon))
	params.Set("toPlace", fmt.Sprintf("%f,%f", req.ToLat, req.ToLon))
	params.Set("date", req.Date)
	params.Set("time", req.Time)
	params.Set("mode", "TRANSIT,WALK")
	params.Set("numItineraries", "3")

	resp, err := http.Get(otpBase + "/plan?" + params.Encode())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result OTPResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

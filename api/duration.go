package api

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"time"
)

type Duration struct {
	Start time.Time     `json:"start"`
	End   time.Time     `json:"end"`
	Value time.Duration `json:"duration"`
}

// JSON implementation for WebFormatter interface.
func (d Duration) JSON() ([]byte, error) {
	return json.Marshal(d)
}

// XML implementation for WebFormatter interface.
func (d Duration) XML() ([]byte, error) {
	return xml.Marshal(d)
}

// Print a representation of Duration to stdout
func (d Duration) Print() {
	// For initial debug purposes...
	fmt.Printf("------------------------------------------------\n")
	fmt.Printf("Day start     : %d\n", d.Start.Day())
	fmt.Printf("Day end       : %d\n", d.End.Day())
	fmt.Printf("Weekday start : %s\n", d.Start.Weekday())
	fmt.Printf("Weekday end   : %s\n", d.End.Weekday())
	fmt.Printf("Mongth start  : %s\n", d.Start.Month())
	fmt.Printf("Mongth end    : %s\n", d.Start.Month())
	fmt.Printf("Time start    : %s\n", d.Start)
	fmt.Printf("Time end      : %s\n", d.End)
	fmt.Printf("Duration      : %s\n", d.Value)
}

func ParseDuration(startParam string, endParam string) (Duration, error) {
	var duration Duration

	startTime, err := time.Parse(time.RFC3339, startParam)
	if err != nil {
		return duration, fmt.Errorf("could not parse 'start' parameter [%s]: %v", startParam, err)
	}

	endTime, err := time.Parse(time.RFC3339, endParam)
	if err != nil {
		return duration, fmt.Errorf("could not parse 'end' parameter [%s]: %v", endParam, err)
	}

	// The end time must come after the start time.
	if endTime.Before(startTime) {
		return duration, fmt.Errorf("end time occurs before start time")
	}

	// Calculate duration from start to end.
	duration.Value = endTime.Sub(startTime)
	duration.Start = startTime
	duration.End = endTime

	return duration, nil
}

func DurationFromHTTPRequest(r *http.Request) (Duration, error) {
	startParam := r.URL.Query()["start"][0]
	endParam := r.URL.Query()["end"][0]
	return ParseDuration(startParam, endParam)
}

// DurationHandleFunc provides an endpoint to that echos back both a start and end timestamp
// in RFC3339 format, after parsing and computing duration.
//
// Example:
// 		curl  "http://localhost:8080/api/duration?start=2015-07-01T07%3A00%3A00Z&end=2015-07-01T12%3A00%3A00Z"
func DurationHandleFunc(w http.ResponseWriter, r *http.Request) {
	InitializeResponse(&w, r) // Required before WriteResponse

	// Calculate duration from start to end.
	duration, err := DurationFromHTTPRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		WriteResponse(APIStandardResponse{http.StatusBadRequest, err.Error()}, &w)
		return
	}

	// Return time formatted as RFC3339 as a sanity check.
	w.WriteHeader(http.StatusOK)
	WriteResponse(duration, &w)
}

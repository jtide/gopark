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

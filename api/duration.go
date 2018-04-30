package api

import (
	"fmt"
	"net/http"
	"time"
)

type Duration struct {
	Start time.Time
	End   time.Time
	Value time.Duration
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

	// The end time must come after the start time
	if endTime.Before(startTime) {
		return duration, fmt.Errorf("end time occurs before start time")
	}

	// Calculate duration from start to end
	duration.Value = endTime.Sub(startTime)
	duration.Start = startTime
	duration.End = endTime

	// For initial debug purposes...
	fmt.Printf("------------------------------------------------\n")
	fmt.Printf("Day start     : %d\n", duration.Start.Day())
	fmt.Printf("Day end       : %d\n", duration.End.Day())
	fmt.Printf("Weekday start : %s\n", duration.Start.Weekday())
	fmt.Printf("Weekday end   : %s\n", duration.End.Weekday())
	fmt.Printf("Mongth start  : %s\n", duration.Start.Month())
	fmt.Printf("Mongth end    : %s\n", duration.Start.Month())
	fmt.Printf("Time start    : %s\n", duration.Start)
	fmt.Printf("Time end      : %s\n", duration.End)
	fmt.Printf("Duration      : %s\n", duration.Value)

	return duration, nil
}

func DurationFromHTTPRequest(r *http.Request) (Duration, error) {
	startParam := r.URL.Query()["start"][0]
	endParam := r.URL.Query()["end"][0]
	return ParseDuration(startParam, endParam)
}

// DurationHandleFunc provides an endpoint to that echos back both a start and end timestamp
// in RFC3339 format, after parsing and computing duration
//
// Example:
// 		curl  "http://localhost:8080/api/duration?start=2015-07-01T07%3A00%3A00Z&end=2015-07-01T12%3A00%3A00Z"
func DurationHandleFunc(w http.ResponseWriter, r *http.Request) {
	InitializeResponse(&w, r)

	// Calculate duration from start to end
	duration, err := DurationFromHTTPRequest(r)
	if err != nil {
		DescribeError(&w, err.Error())
		return
	}

	// Return time formatted as RFC3339 as a sanity check
	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "text/plain")
	fmt.Fprintf(w, "start    : %s\n", duration.Start)
	fmt.Fprintf(w, "end      : %s\n", duration.End)
	fmt.Fprintf(w, "duration : %s\n", duration.Value)
}

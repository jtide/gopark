package main

import (
	"net/http"
	"fmt"
	"os"
	"time"
)

func main() {
	http.HandleFunc("/api/duration", duration)
	http.ListenAndServe(port(), nil)
}

func port() string {
	port := os.Getenv("GOPARK_PORT")
	if len(port) == 0 {
		port = "8080"
	}
	return ":" + port
}

// duration provides an endpoint to that echos back both a start and end timestamp
// in RFC3339 format, after parsing and computing duration
//
// Example:
// 		curl  "http://localhost:8080/api/duration?start=2015-07-01T07%3A00%3A00Z&end=2015-07-01T12%3A00%3A00Z"
func duration(w http.ResponseWriter, r *http.Request)  {
	startParam := r.URL.Query()["start"][0]
	endParam := r.URL.Query()["end"][0]

	start, err := time.Parse(time.RFC3339, startParam)
	if err != nil {
		fmt.Printf("Error parsing startParam %s\n", startParam)
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Add("Content-Type", "text/plain")
		fmt.Fprintf(w, "Invalid start time: %s\n", start)
		return
	}

	end, err := time.Parse(time.RFC3339, endParam)
	if err != nil {
		fmt.Printf("Error parsing endParam %s\n", endParam)
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Add("Content-Type", "text/plain")
		fmt.Fprintf(w, "Invalid end time: %s\n", start)
		return
	}

	// The end time must come after the start time
	if end.Before(start) {
		fmt.Printf("Error computing duration: %s > %s\n", start, end)
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Add("Content-Type", "text/plain")
		fmt.Fprintf(w, "Invalid duration: End time occurs before start time\n")
		return
	}

	// Calculate duration from start to end
	dur := end.Sub(start)

	// For initial debug purposes...
	fmt.Printf("------------------------------------------------\n")
	fmt.Printf("Day start     : %d\n", start.Day())
	fmt.Printf("Day end       : %d\n", end.Day())
	fmt.Printf("Weekday start : %s\n", start.Weekday())
	fmt.Printf("Weekday end   : %s\n", end.Weekday())
	fmt.Printf("Mongth start  : %s\n", start.Month())
	fmt.Printf("Mongth end    : %s\n", start.Month())
	fmt.Printf("Time start    : %s\n", start)
	fmt.Printf("Time end      : %s\n", end)
	fmt.Printf("Duration      : %s\n", dur)

	// Return time formatted as RFC3339 as a sanity check
	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "text/plain")
	fmt.Fprintf(w, "start    : %s\n", start)
	fmt.Fprintf(w, "end      : %s\n", end)
	fmt.Fprintf(w, "duration : %s\n", dur)
}

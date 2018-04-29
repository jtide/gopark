package main

import (
	"net/http"
	"fmt"
	"os"
	"time"
	"encoding/json"
	"encoding/xml"
	"strings"
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

// WebFormatter can produce either XML or JSON representations of itself
type WebFormatter interface {
	Json() []byte
	Xml() []byte
}

// ApiError implements WebFormatter interface
type ApiError struct {
	Description string
}

func (e ApiError) Json() []byte {
	response, err := json.Marshal(e)
	if err != nil {
		panic(err)
	}
	return response
}

func (e ApiError) Xml() []byte {
	response, err := xml.Marshal(e)
	if err != nil {
		panic(err)
	}
	return response
}

func respondWithJson(f WebFormatter, w *http.ResponseWriter) {
	(*w).Header().Add("Content-Type", "application/json; charset-utf-8")
	(*w).Write(f.Json())
}

func respondWithXml(f WebFormatter, w *http.ResponseWriter) {
	(*w).Header().Add("Content-Type", "application/xml; charset-utf-8")
	(*w).Write(f.Xml())
}

// respond is responsible for writing the response payload in either XML or JSON format, based on
// what the client specifies is accepted in the HTTP headers. If not otherwise specified, JSON is
// used by default.
func respond(f WebFormatter, w *http.ResponseWriter, r *http.Request) {
	encoding := r.Header.Get("Accept")
	switch {
	case strings.Contains(encoding, "application/json"):
		respondWithJson(f, w)
	case strings.Contains(encoding, "application/xml"):
		respondWithXml(f, w)
	default:
		respondWithJson(f, w)
	}
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
		e := ApiError{ fmt.Sprintf("Error parsing 'start' parameter: [%s]", startParam)}
		w.WriteHeader(http.StatusBadRequest)
		respond(e, &w, r)
		return
	}

	end, err := time.Parse(time.RFC3339, endParam)
	if err != nil {
		e := ApiError{ fmt.Sprintf("Error parsing 'end' parameter: [%s]", endParam)}
		w.WriteHeader(http.StatusBadRequest)
		respond(e, &w, r)
		return
	}

	// The end time must come after the start time
	if end.Before(start) {
		e := ApiError{ "Invalid duration: End time occurs before start time"}
		w.WriteHeader(http.StatusBadRequest)
		respond(e, &w, r)
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

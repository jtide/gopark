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
// the response Content-Type HTTP header. If not otherwise specified, JSON is used by default.
func respond(f WebFormatter, w *http.ResponseWriter) {
	encoding := (*w).Header().Get("Content-Type")
	switch {
	case strings.Contains(encoding, "json"):
		respondWithJson(f, w)
	case strings.Contains(encoding, "xml"):
		respondWithXml(f, w)
	default:
		respondWithJson(f, w)
	}
}

// setResponseFormat sets the format of the response based on the "Accept" headers of the HTTP request.
// The format will be either JSON or XML.  If the client accepts either, then JSON is preferred.
func setResponseFormat(w *http.ResponseWriter, r *http.Request) {
	encoding := r.Header.Get("Accept")
	switch {
	case strings.Contains(encoding, "application/json"):
		(*w).Header().Add("Content-Type", "application/json; charset-utf-8")
	case strings.Contains(encoding, "application/xml"):
		(*w).Header().Add("Content-Type", "application/xml; charset-utf-8")
	default:
		(*w).Header().Add("Content-Type", "application/json; charset-utf-8")
	}
}

// describeError automatically populates an http response with an error message, appropriately
// formatted in either JSON or XML.
func describeError(w *http.ResponseWriter, description string) {
	e := ApiError{ description }
	(*w).WriteHeader(http.StatusBadRequest)
	respond(e, w)
}

// duration provides an endpoint to that echos back both a start and end timestamp
// in RFC3339 format, after parsing and computing duration
//
// Example:
// 		curl  "http://localhost:8080/api/duration?start=2015-07-01T07%3A00%3A00Z&end=2015-07-01T12%3A00%3A00Z"
func duration(w http.ResponseWriter, r *http.Request)  {
	setResponseFormat(&w, r)

	startParam := r.URL.Query()["start"][0]
	endParam := r.URL.Query()["end"][0]


	start, err := time.Parse(time.RFC3339, startParam)
	if err != nil {
		describeError(&w, fmt.Sprintf("Error parsing 'start' parameter: [%s]", startParam))
		return
	}

	end, err := time.Parse(time.RFC3339, endParam)
	if err != nil {
		describeError(&w, fmt.Sprintf("Error parsing 'end' parameter: [%s]", endParam))
		return
	}

	// The end time must come after the start time
	if end.Before(start) {
		describeError(&w, "Invalid duration: End time occurs before start time")
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

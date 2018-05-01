package api

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
	"time"
)

type Rate struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
	Price uint      `json:"price"`
}

// JSON implementation for WebFormatter interface.
func (r Rate) JSON() ([]byte, error) {
	return json.Marshal(r)
}

// XML implementation for WebFormatter interface.
func (r Rate) XML() ([]byte, error) {
	return xml.Marshal(r)
}

func lookupRate(d Duration) uint {
	// TODO: implement lookupRate
	return uint(0)
}

// RateHandleFunc provides an endpoint to that echos back both a start and end timestamp
// in RFC3339 format, after parsing and computing duration.
//
// Example:
// 		curl  "http://localhost:8080/api/duration?start=2015-07-01T07%3A00%3A00Z&end=2015-07-01T12%3A00%3A00Z"
func RateHandleFunc(w http.ResponseWriter, r *http.Request) {
	InitializeResponse(&w, r) // Required before WriteResponse

	// Calculate duration from start to end
	duration, err := DurationFromHTTPRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		WriteResponse(Error{err.Error()}, &w)
		return
	}

	// Lookup the Rate
	rate := Rate{Start: duration.Start, End: duration.End, Price: lookupRate(duration)}

	// Respond with price
	w.WriteHeader(http.StatusOK)
	err = WriteResponse(rate, &w)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		WriteResponse(Error{err.Error()}, &w)
		return
	}
}

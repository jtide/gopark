package api

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
)

type Error struct {
	Description string
}

// JSON implementation for WebFormatter interface
func (e Error) JSON() []byte {
	response, err := json.Marshal(e)
	if err != nil {
		panic(err)
	}
	return response
}

// XML implementation for WebFormatter interface
func (e Error) XML() []byte {
	response, err := xml.Marshal(e)
	if err != nil {
		panic(err)
	}
	return response
}

// DescribeError populates an http response with an error message, appropriately
// formatted in either JSON or XML.
func DescribeError(w *http.ResponseWriter, description string) {
	e := Error{description}
	(*w).WriteHeader(http.StatusBadRequest)
	Respond(e, w)
}

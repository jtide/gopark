package api

import (
	"encoding/json"
	"encoding/xml"
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

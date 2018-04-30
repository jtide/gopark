package api

import (
	"encoding/json"
	"encoding/xml"
)

type Error struct {
	Description string
}

// JSON implementation for WebFormatter interface
func (e Error) JSON() ([]byte, error) {
	return json.Marshal(e)
}

// XML implementation for WebFormatter interface
func (e Error) XML() ([]byte, error) {
	return xml.Marshal(e)
}

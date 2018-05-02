package api

import (
	"encoding/json"
	"encoding/xml"
)

type APIError struct {
	Description string `json:"error"`
}

// JSON implementation for WebFormatter interface.
func (e APIError) JSON() ([]byte, error) {
	return json.Marshal(e)
}

// XML implementation for WebFormatter interface.
func (e APIError) XML() ([]byte, error) {
	return xml.Marshal(e)
}

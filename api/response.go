package api

import (
	"net/http"
	"strings"
	"encoding/json"
	"encoding/xml"
)

// WebFormatter implementors can produce either XML or JSON representations themselves.
type WebFormatter interface {
	JSON() ([]byte, error)
	XML() ([]byte, error)
}

// APIStandardResponse provides a vehicle for generic success/error responses
type APIStandardResponse struct {
	Status      uint   `json:"status"`
	Description string `json:"desc"`
}

// JSON implementation for WebFormatter interface.
func (e APIStandardResponse) JSON() ([]byte, error) {
	return json.Marshal(e)
}

// XML implementation for WebFormatter interface.
func (e APIStandardResponse) XML() ([]byte, error) {
	return xml.Marshal(e)
}

func respondWithJSON(f WebFormatter, w *http.ResponseWriter) error {
	response, err := f.JSON()
	if err != nil {
		return err
	}
	(*w).Header().Add("Content-Type", "application/json; charset-utf-8")
	(*w).Write(response)
	return nil
}

func respondWithXML(f WebFormatter, w *http.ResponseWriter) error {
	response, err := f.XML()
	if err != nil {
		return err
	}
	(*w).Header().Add("Content-Type", "application/xml; charset-utf-8")
	(*w).Write(response)
	return nil
}

// WriteResponse is responsible for writing the response payload in either XML or JSON format, based on
// the response Content-Type HTTP header. If not otherwise specified, JSON is used by default.
func WriteResponse(f WebFormatter, w *http.ResponseWriter) error {
	encoding := (*w).Header().Get("Content-Type")
	switch {
	case strings.Contains(encoding, "json"):
		return respondWithJSON(f, w)
	case strings.Contains(encoding, "xml"):
		return respondWithXML(f, w)
	default:
		return respondWithJSON(f, w)
	}
}

// InitializeResponse sets the format of the response based on the "Accept" headers of the HTTP request.
// InitializeResponse must be called before WriteResponse in order to ensure proper format of the response.
// The format will be either JSON or XML.  If the client accepts either, then JSON is preferred.
func InitializeResponse(w *http.ResponseWriter, r *http.Request) {
	encoding := r.Header.Get("Accept")
	switch {
	case strings.Contains(encoding, "json"):
		(*w).Header().Add("Content-Type", "application/json; charset-utf-8")
	case strings.Contains(encoding, "xml"):
		(*w).Header().Add("Content-Type", "application/xml; charset-utf-8")
	default:
		(*w).Header().Add("Content-Type", "application/json; charset-utf-8")
	}
}

package api

import (
	"net/http"
	"strings"
)

// WebFormatter can produce either XML or JSON representations of itself
type WebFormatter interface {
	JSON() ([]byte, error)
	XML() ([]byte, error)
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
// The format will be either JSON or XML.  If the client accepts either, then JSON is preferred.
func InitializeResponse(w *http.ResponseWriter, r *http.Request) {
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

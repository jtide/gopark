package api

import (
	"net/http"
	"strings"
)

// WebFormatter can produce either XML or JSON representations of itself
type WebFormatter interface {
	JSON() []byte
	XML() []byte
}

func respondWithJSON(f WebFormatter, w *http.ResponseWriter) {
	(*w).Header().Add("Content-Type", "application/json; charset-utf-8")
	(*w).Write(f.JSON())
}

func respondWithXML(f WebFormatter, w *http.ResponseWriter) {
	(*w).Header().Add("Content-Type", "application/xml; charset-utf-8")
	(*w).Write(f.XML())
}

// WriteResponse is responsible for writing the response payload in either XML or JSON format, based on
// the response Content-Type HTTP header. If not otherwise specified, JSON is used by default.
func WriteResponse(f WebFormatter, w *http.ResponseWriter) {
	encoding := (*w).Header().Get("Content-Type")
	switch {
	case strings.Contains(encoding, "json"):
		respondWithJSON(f, w)
	case strings.Contains(encoding, "xml"):
		respondWithXML(f, w)
	default:
		respondWithJSON(f, w)
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

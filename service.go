package main

import (
	"github.com/jtide/gopark/api"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/api/duration", api.DurationHandleFunc)
	http.HandleFunc("/api/rate", api.RateHandleFunc)
	http.ListenAndServe(port(), nil)
}

func port() string {
	port := os.Getenv("GOPARK_PORT")
	if len(port) == 0 {
		port = "8080"
	}
	return ":" + port
}

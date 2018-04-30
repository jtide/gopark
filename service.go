package main

import (
	"net/http"
	"os"
	"github.com/jtide/gopark/api"
)

func main() {
	http.HandleFunc("/api/duration", api.DurationHandleFunc)
	http.ListenAndServe(port(), nil)
}

func port() string {
	port := os.Getenv("GOPARK_PORT")
	if len(port) == 0 {
		port = "8080"
	}
	return ":" + port
}

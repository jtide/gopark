package main

import (
	"net/http"
	"fmt"
	"os"
)

func main() {
	http.HandleFunc("/", index)
	http.ListenAndServe(port(), nil)

}

func port() string {
	port := os.Getenv("GOPARK_PORT")
	if len(port) == 0 {
		port = "8080"
	}
	return ":" + port
}

func index(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Welcome to My New Parking App!\n")
}

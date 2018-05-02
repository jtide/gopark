package main

import (
	"flag"
	"fmt"
	"github.com/jtide/gopark/api"
	"io/ioutil"
	"net/http"
	"os"
)

func main() {
	err := api.ReplaceRates(jsonRateConfig())
	if err != nil {
		// Failure to apply the initial rate configuration is
		// one of the very few cases where a panic is warranted.
		panic(err)
	}
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

func jsonRateConfig() []byte {
	configFile := flag.String("config", "", "Absolute path to JSON rate configuration file.")
	flag.Parse()

	if *configFile == "" {
		return api.JSONDefaultRateConfig
	}

	fmt.Println("Using rates configuration file:", *configFile)
	configJSON, err := ioutil.ReadFile(*configFile)
	if err != nil {
		// Panic if configuration file cannot be read
		panic(err)
	}

	return configJSON
}

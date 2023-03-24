package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/callsign", getCallSignInfo)

	mux.HandleFunc("/weather", getWeatherInfo)

	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}

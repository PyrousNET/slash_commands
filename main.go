package main

import (
	"log"
	"net/http"
)

type (
	Callsign struct {
		Call string
	}
	Response struct {
		ResponseType string `json:"response_type"`
		Text         string `json:"text"`
	}
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/callsign", getCallSignInfo)

	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}

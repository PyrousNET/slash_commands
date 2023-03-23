package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/callsign", getCallSignInfo)

	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}

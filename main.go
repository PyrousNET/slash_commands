package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	fmt.Println("Starting...")
	mux := http.NewServeMux()

	fmt.Println("Setting up Callsign route.")
	mux.HandleFunc("/callsign", getCallSignInfo)
	fmt.Println("Setting up Book route.")
	mux.HandleFunc("/book", getBookInfo)
	fmt.Println("Setting up Weather route.")
	mux.HandleFunc("/weather", getWeatherInfo)

	fmt.Println("Listening for connections on port 4000")
	err := http.ListenAndServe(":4000", mux)

	log.Fatal(err)
}

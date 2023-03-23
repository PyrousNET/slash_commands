package main

import (
	"encoding/json"
	"fmt"
	"github.com/pyrousnet/slash_commands/Callook"
	"github.com/pyrousnet/slash_commands/HamDb"
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

func getCallSignInfo(w http.ResponseWriter, r *http.Request) {
	text := r.URL.Query().Get("text")
	var c Callsign = Callsign{Call: text}

	err, hCS := HamDb.PullFromHamDb(c.Call)
	err = fmt.Errorf("test", 500)
	if err != nil {
		err, hCS = Callook.PullFromCallook(c.Call)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}

	mmrstr := hCS.GetResponseString()

	rs := Response{
		ResponseType: "in_channel",
		Text:         string(mmrstr),
	}

	b, err := json.Marshal(rs)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return
	}

	w.Write([]byte(b))
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/callsign", getCallSignInfo)

	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}

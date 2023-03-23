package main

import (
	"encoding/json"
	"fmt"
	"github.com/pyrousnet/slash_commands/Callook"
	"github.com/pyrousnet/slash_commands/HamDb"
	"github.com/pyrousnet/slash_commands/MatterMost"
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

	err, mmr := HamDb.PullFromHamDb(c.Call)
	err = fmt.Errorf("test", 500)
	if err != nil {
		err, mmr = Callook.PullFromCallook(c.Call)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}

	mmrstr := getMatterMostStringFromMMResponse(mmr)

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

func getMatterMostStringFromMMResponse(r *MatterMost.Response) string {
	return "| Data | Value |\n| :------ | :-------|\n| Callsign | " + r.CallSign +
		" |\n| Name | " + r.Name +
		" |\n| City | " + r.City +
		" |\n| Last3 | " + r.Last3 +
		" |\n| Class | " + r.Class +
		" |\n| Status | " + r.Status + " |"
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/callsign", getCallSignInfo)

	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}

package main

import (
	"encoding/json"
	"fmt"
	"github.com/pyrousnet/slash_commands/Callook"
	"github.com/pyrousnet/slash_commands/Color"
	"github.com/pyrousnet/slash_commands/HamDb"
	"github.com/pyrousnet/slash_commands/MatterMost"
	"net/http"
)

type Callsign struct {
	Call string
}

func getCallSignInfo(w http.ResponseWriter, r *http.Request) {
	text := r.URL.Query().Get("text")
	var c Callsign = Callsign{Call: text}

	formattedText := Color.Reset + Color.Cyan + text + Color.Reset

	fmt.Printf(Color.Green + "Incomming callsign request for: " + formattedText + "\n")

	err, hCS := HamDb.PullFromHamDb(c.Call)
	if err != nil {
		err, hCS = Callook.PullFromCallook(c.Call)
		if err != nil {
			formattedText = Color.Red + "Error: " + err.Error() + Color.Reset
			fmt.Printf(formattedText + "\n")
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}

	mmrstr := hCS.GetResponseString()

	rs := MatterMost.Response{
		ResponseType: "in_channel",
		Text:         string(mmrstr),
	}

	b, err := json.Marshal(rs)
	if err != nil {
		formattedText = Color.Red + "Error: " + err.Error() + Color.Reset
		fmt.Printf(formattedText + "\n")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	fmt.Println("Sending response")
	w.Write([]byte(b))
}

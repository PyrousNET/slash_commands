package main

import (
	"encoding/json"
	"fmt"
	"github.com/pyrousnet/slash_commands/Callook"
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

	fmt.Printf("\033[32mIncomming callsign request for:\033[0m\033[36m " + text + "\033[0m\n")

	err, hCS := HamDb.PullFromHamDb(c.Call)
	if err != nil {
		err, hCS = Callook.PullFromCallook(c.Call)
		if err != nil {
			fmt.Printf("\033[31mError: " + err.Error() + "\033[0m\n")
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
		fmt.Printf("\033[31mError: " + err.Error() + "\033[0m\n")
		fmt.Printf("Error: %s", err)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte(b))
}

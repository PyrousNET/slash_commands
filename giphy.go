package main

import (
	"encoding/json"
	"fmt"
	"github.com/pyrousnet/slash_commands/Color"
	"github.com/pyrousnet/slash_commands/Giphy" // Import the Giphy package
	"github.com/pyrousnet/slash_commands/MatterMost"
	"log"
	"net/http"
	"os"
)

var giphyStates = make(map[string]Giphy.GiphyState)

func giphyCommand(w http.ResponseWriter, r *http.Request) {
	text := r.URL.Query().Get("text")
	formattedText := Color.Reset + Color.Cyan + text + Color.Reset
	fmt.Printf(Color.Green+"Incomming giphy request for: %s\n", formattedText+Color.Reset)

	apiKey := os.Getenv("GIPHY_API_KEY")
	g := Giphy.Setup(apiKey) // Replace YOUR_GIPHY_API_KEY with your actual Giphy API key
	err, giphyResponse := g.PullFromGiphy(text)
	if err != nil {
		http.Error(w, "Error calling GIPHY API", http.StatusInternalServerError)
		fmt.Println(Color.Red + "GIPHY API Error: " + err.Error() + Color.Reset)
		return
	}

	userName := r.URL.Query().Get("user_name")
	giphyStates[text] = Giphy.GiphyState{Results: *giphyResponse, CurrentIndex: 0, SearchTerm: text, User: userName}

	mmResponse := Giphy.CreateMatterMostResponse(*giphyResponse, 0, userName, text)
	responseJson, err := json.Marshal(mmResponse)
	if err != nil {
		log.Println(Color.Red + "Error marshalling JSON:" + err.Error() + Color.Reset)
		http.Error(w, "Error marshalling JSON", http.StatusInternalServerError)
		return
	}
	fmt.Println(string(responseJson))
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(mmResponse)
	if err != nil {
		log.Println(Color.Red + "Error encoding JSON:" + err.Error() + Color.Reset)
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		return
	}
	fmt.Println(Color.Green + "GIPHY request successful" + Color.Reset)
}

func sendGiphyUpdate(w http.ResponseWriter, imageUrl string) {
	type (
		Attachments []struct {
			ImageURL string `json:"image_url"`
		}
		Props struct {
			Attachments Attachments `json:"attachments"`
		}
		Update struct {
			Message string `json:"message"`
			Props   Props  `json:"props"`
		}
		Response struct {
			Update Update `json:"update"`
		}
	)

	mmResponse := Response{
		Update: Update{
			Message: imageUrl,
			Props: Props{
				Attachments: Attachments{
					{
						ImageURL: imageUrl,
					},
				},
			},
		},
	}
	err := json.NewEncoder(w).Encode(mmResponse)
	if err != nil {
		log.Println(Color.Red + "Error encoding JSON:" + err.Error() + Color.Reset)
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		return
	}
	fmt.Println(Color.Green + "GIPHY request successful" + Color.Reset)
}

func sendGiphyPreview(w http.ResponseWriter, key string) {
	state := giphyStates[key]
	mmResponse := Giphy.CreateMatterMostResponse(state.Results, state.CurrentIndex, state.User, state.SearchTerm)
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(mmResponse)
	if err != nil {
		log.Println(Color.Red + "Error encoding JSON:" + err.Error() + Color.Reset)
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		return
	}
	fmt.Println(Color.Green + "GIPHY request successful" + Color.Reset)
}

func giphyPrevious(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("text")
	state := giphyStates[key]
	if state.CurrentIndex > 0 {
		state.CurrentIndex--
		giphyStates[key] = state
	}
	sendGiphyPreview(w, key)
}

func giphyNext(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("text")
	state := giphyStates[key]
	if state.CurrentIndex < len(state.Results.Data)-1 {
		state.CurrentIndex++
		giphyStates[key] = state
	}
	//sendGiphyPreview(w, key)
	sendGiphyUpdate(w, state.Results.Data[state.CurrentIndex].Images.Original.URL)
}

func giphySelect(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("text")
	state := giphyStates[key]
	if state.CurrentIndex < 0 {
		state.CurrentIndex = 0
	}
	if state.CurrentIndex >= len(state.Results.Data) {
		state.CurrentIndex = len(state.Results.Data) - 1
	}
	gif := state.Results.Data[state.CurrentIndex]
	originalURL := gif.Images.Original.URL
	userName := r.URL.Query().Get("user_name")

	mmResponse := MatterMost.Response{
		ResponseType: "in_channel",
		Text:         fmt.Sprintf("%s Posted GIF: ![%s](%s)", userName, key, originalURL),
	}
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(mmResponse)
	if err != nil {
		log.Println(Color.Red + "Error encoding JSON:" + err.Error() + Color.Reset)
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		return
	}
	delete(giphyStates, key)
}

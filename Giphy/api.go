package Giphy

import (
	"encoding/json"
	"fmt"
	"github.com/pyrousnet/slash_commands/MatterMost"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

type Giphy struct {
	apiKey string `json:"apiKey"`
	host   string `json:"host"`
	limit  int    `json:"limit"`
}

func Setup(key string) *Giphy {
	g := Giphy{
		apiKey: key,
		host:   "https://api.giphy.com",
		limit:  5,
	}
	return &g
}

func (g *Giphy) PullFromGiphy(searchTerm string) (error, *Response) {
	var r Response
	searchURL := g.host + "/v1/gifs/search?api_key=" + g.apiKey + "&q=" + url.QueryEscape(searchTerm) + "&limit=" + strconv.Itoa(g.limit)

	resp, err := http.Get(searchURL)
	if err != nil {
		return err, nil
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err, nil
	}

	err = json.Unmarshal(body, &r)
	if err != nil {
		return err, nil
	}
	r.Query = searchTerm
	return nil, &r
}

func CreateMatterMostResponse(r Response, currentIndex int, user string, searchTerm string) MatterMost.Response {
	if len(r.Data) == 0 {
		return MatterMost.Response{
			ResponseType: "in_channel",
			Text:         "No GIFs found for that search.",
		}
	}
	gif := r.Data[currentIndex]
	previewURL := gif.Images.Preview.URL

	// Construct the text to include the preview and buttons
	responseText := fmt.Sprintf("Previewing GIF... (%d/%d)\n", currentIndex+1,
		len(r.Data),
	)
	/*
		responseText += fmt.Sprintf("[Previous](/giphy/previous?key=%s) | ", searchTerm)
		responseText += fmt.Sprintf("[Next](/giphy/next?key=%s) | ", searchTerm)
		//responseText += fmt.Sprintf("[Select](/giphy/select?key=%s)", searchTerm)
		// I would like this to launch a slash command in mattermost to post the gif
		// to the channel, but I'm not sure how to do that yet and not open the gif in a new tab.
		// This is done using the MatterMost slash command API, but I'm not sure how to do that yet.
		responseText
	*/

	attachments := []MatterMost.Attachment{
		{
			ImageUrl: previewURL,
			Actions: []MatterMost.Action{
				{
					Id:   "gif-previous",
					Name: "◀️ Previous",
					Integration: MatterMost.Integration{
						URL: "/giphy/previous?text=" + searchTerm,
						Context: map[string]string{
							"current_index": strconv.Itoa(currentIndex),
							"direction":     "previous",
						},
					},
				},
				{
					Id:   "gif-next",
					Name: "▶️ Next",
					Integration: MatterMost.Integration{
						URL: "http://host.docker.internal:4000/giphy/next?text=" + searchTerm,
						Context: map[string]string{
							"current_index": strconv.Itoa(currentIndex),
							"direction":     "next",
						},
					},
				},
				{
					Id:   "gif-select",
					Name: "✅ Select",
					Integration: MatterMost.Integration{
						URL: "/giphy/select?text=" + searchTerm,
						Context: map[string]string{
							"current_index": strconv.Itoa(currentIndex),
						},
					},
				},
			},
		},
	}

	return MatterMost.Response{
		ResponseType: "ephemeral",
		Text:         responseText,
		Attachments:  attachments,
	}
}

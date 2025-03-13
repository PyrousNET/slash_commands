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
	return nil, &r
}

func CreateMatterMostResponse(r Response, currentIndex int, searchTerm string) MatterMost.Response {
	if len(r.Data) == 0 {
		return MatterMost.Response{
			ResponseType: "in_channel",
			Text:         "No GIFs found for that search.",
		}
	}
	gif := r.Data[currentIndex]
	previewURL := gif.Images.Preview.URL

	// Construct the text to include the preview and buttons
	responseText := fmt.Sprintf("Previewing GIF... (%d/%d)\n%s\n", currentIndex+1, len(r.Data), previewURL)
	responseText += fmt.Sprintf("[Previous](/giphy/previous?key=%s) | ", searchTerm)
	responseText += fmt.Sprintf("[Next](/giphy/next?key=%s) | ", searchTerm)
	responseText += fmt.Sprintf("[Select](/giphy/select?key=%s)", searchTerm)

	return MatterMost.Response{
		ResponseType: "in_channel",
		Text:         responseText,
	}
}

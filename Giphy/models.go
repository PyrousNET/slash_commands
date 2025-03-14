package Giphy

import (
	_ "github.com/pyrousnet/slash_commands/MatterMost"
)

type Response struct {
	Query string
	Data  []struct {
		Id     string `json:"id"`
		Images struct {
			Preview struct {
				URL string `json:"url"`
			} `json:"downsized"`
			Original struct {
				URL string `json:"url"`
			} `json:"original"`
		} `json:"images"`
	} `json:"data"`
}

type GiphyState struct {
	Results      Response
	CurrentIndex int
	SearchTerm   string
	User         string
}

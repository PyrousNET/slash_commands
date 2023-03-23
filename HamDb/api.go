package HamDb

import (
	"encoding/json"
	"github.com/pyrousnet/slash_commands/MatterMost"
	"io/ioutil"
	"net/http"
)

type (
	Response struct {
		Hamdb struct {
			Version  string `json:"version"`
			Callsign struct {
				Call    string `json:"call"`
				Class   string `json:"class"`
				Expires string `json:"expires"`
				Status  string `json:"status"`
				Grid    string `json:"grid"`
				Lat     string `json:"lat"`
				Lon     string `json:"lon"`
				Fname   string `json:"fname"`
				Mi      string `json:"mi"`
				Name    string `json:"name"`
				Suffix  string `json:"suffix"`
				Addr1   string `json:"addr1"`
				Addr2   string `json:"addr2"`
				State   string `json:"state"`
				Zip     string `json:"zip"`
				Country string `json:"country"`
			} `json:"callsign"`
			Messages struct {
				Status string `json:"status"`
			} `json:"messages"`
		} `json:"hamdb"`
	}
)

func PullFromHamDb(callsign string) (error, *MatterMost.Response) {
	var r Response
	var mmr MatterMost.Response

	resp, err := http.Get("http://api.hamdb.org/v1/" + callsign + "/json")
	if err != nil {
		return err, nil
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err, nil
	}

	err = json.Unmarshal([]byte(body), &r)
	if err != nil {
		return err, nil
	}

	mmr = createMatterMostResponse(r)

	return nil, &mmr
}

func createMatterMostResponse(r Response) MatterMost.Response {
	return MatterMost.Response{
		CallSign: r.Hamdb.Callsign.Call,
		Name:     r.Hamdb.Callsign.Fname + " " + r.Hamdb.Callsign.Mi + " " + r.Hamdb.Callsign.Name,
		City:     r.Hamdb.Callsign.Addr2 + ", " + r.Hamdb.Callsign.State + " " + r.Hamdb.Callsign.Zip,
		Last3:    r.Hamdb.Callsign.Call[len(r.Hamdb.Callsign.Call)-3:],
		Class:    getClass(r.Hamdb.Callsign.Class),
		Status:   getStatus(r.Hamdb.Callsign.Status),
	}
}

func getClass(c string) string {
	switch c {
	case "T":
		return "Technician"
	case "G":
		return "General"
	case "E":
		return "Extra"
	case "N":
		return "Novice"
	case "A":
		return "Advanced"
	default:
		return "Unknown"
	}
}

func getStatus(s string) string {
	switch s {
	case "A":
		return "Active"
	case "E":
		return "Expired"
	case "C":
		return "Cancelled"
	default:
		return "Unknown"
	}
}

package Callook

import (
	"encoding/json"
	"github.com/pyrousnet/slash_commands/MatterMost"
	"io/ioutil"
	"net/http"
)

type Response struct {
	Status  string `json:"status"`
	Type    string `json:"type"`
	Current struct {
		Callsign  string `json:"callsign"`
		OperClass string `json:"operClass"`
	} `json:"current"`
	Previous struct {
		Callsign  string `json:"callsign"`
		OperClass string `json:"operClass"`
	} `json:"previous"`
	Trustee struct {
		Callsign string `json:"callsign"`
		Name     string `json:"name"`
	} `json:"trustee"`
	Name    string `json:"name"`
	Address struct {
		Line1 string `json:"line1"`
		Line2 string `json:"line2"`
		Attn  string `json:"attn"`
	} `json:"address"`
	Location struct {
		Latitude   string `json:"latitude"`
		Longitude  string `json:"longitude"`
		Gridsquare string `json:"gridsquare"`
	} `json:"location"`
	OtherInfo struct {
		GrantDate      string `json:"grantDate"`
		ExpiryDate     string `json:"expiryDate"`
		LastActionDate string `json:"lastActionDate"`
		Frn            string `json:"frn"`
		UlsURL         string `json:"ulsUrl"`
	} `json:"otherInfo"`
}

func PullFromCallook(callsign string) (error, *MatterMost.HamCallSign) {
	var r Response
	var mmr MatterMost.HamCallSign

	resp, err := http.Get("http://callook.info/" + callsign + "/json")
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

func createMatterMostResponse(r Response) MatterMost.HamCallSign {
	return MatterMost.HamCallSign{
		CallSign: r.Current.Callsign,
		Name:     r.Name,
		City:     r.Address.Line2,
		Last3:    r.Current.Callsign[len(r.Current.Callsign)-3:],
		Class:    r.Current.OperClass,
		Status:   r.Status,
	}
}

package main

import (
	"encoding/json"
	"fmt"
	"github.com/pyrousnet/slash_commands/MatterMost"
	"io/ioutil"
	"net/http"
)

type (
	Forcast struct {
		Context  []interface{} `json:"@context"`
		Type     string        `json:"type"`
		Geometry struct {
			Type        string        `json:"type"`
			Coordinates [][][]float64 `json:"coordinates"`
		} `json:"geometry"`
		Properties struct {
			Updated           string `json:"updated"`
			Units             string `json:"units"`
			ForecastGenerator string `json:"forecastGenerator"`
			GeneratedAt       string `json:"generatedAt"`
			UpdateTime        string `json:"updateTime"`
			ValidTimes        string `json:"validTimes"`
			Elevation         struct {
				UnitCode string  `json:"unitCode"`
				Value    float64 `json:"value"`
			} `json:"elevation"`
			Periods []struct {
				Number                     int    `json:"number"`
				Name                       string `json:"name"`
				StartTime                  string `json:"startTime"`
				EndTime                    string `json:"endTime"`
				IsDaytime                  bool   `json:"isDaytime"`
				Temperature                int    `json:"temperature"`
				TemperatureUnit            string `json:"temperatureUnit"`
				TemperatureTrend           string `json:"temperatureTrend"`
				ProbabilityOfPrecipitation struct {
					UnitCode string `json:"unitCode"`
					Value    int    `json:"value"`
				} `json:"probabilityOfPrecipitation"`
				Dewpoint struct {
					UnitCode string  `json:"unitCode"`
					Value    float64 `json:"value"`
				} `json:"dewpoint"`
				RelativeHumidity struct {
					UnitCode string `json:"unitCode"`
					Value    int    `json:"value"`
				} `json:"relativeHumidity"`
				WindSpeed        string `json:"windSpeed"`
				WindDirection    string `json:"windDirection"`
				Icon             string `json:"icon"`
				ShortForecast    string `json:"shortForecast"`
				DetailedForecast string `json:"detailedForecast"`
			} `json:"periods"`
		} `json:"properties"`
	}
	Weather struct {
		Context  []interface{} `json:"@context"`
		ID       string        `json:"id"`
		Type     string        `json:"type"`
		Geometry struct {
			Type        string    `json:"type"`
			Coordinates []float64 `json:"coordinates"`
		} `json:"geometry"`
		Properties struct {
			ID                  string `json:"@id"`
			Type                string `json:"@type"`
			Cwa                 string `json:"cwa"`
			ForecastOffice      string `json:"forecastOffice"`
			GridID              string `json:"gridId"`
			GridX               int    `json:"gridX"`
			GridY               int    `json:"gridY"`
			Forecast            string `json:"forecast"`
			ForecastHourly      string `json:"forecastHourly"`
			ForecastGridData    string `json:"forecastGridData"`
			ObservationStations string `json:"observationStations"`
			RelativeLocation    struct {
				Type     string `json:"type"`
				Geometry struct {
					Type        string    `json:"type"`
					Coordinates []float64 `json:"coordinates"`
				} `json:"geometry"`
				Properties struct {
					City     string `json:"city"`
					State    string `json:"state"`
					Distance struct {
						UnitCode string  `json:"unitCode"`
						Value    float64 `json:"value"`
					} `json:"distance"`
					Bearing struct {
						UnitCode string `json:"unitCode"`
						Value    int    `json:"value"`
					} `json:"bearing"`
				} `json:"properties"`
			} `json:"relativeLocation"`
			ForecastZone    string `json:"forecastZone"`
			County          string `json:"county"`
			FireWeatherZone string `json:"fireWeatherZone"`
			TimeZone        string `json:"timeZone"`
			RadarStation    string `json:"radarStation"`
		} `json:"properties"`
	}
)

func getWeatherInfo(w http.ResponseWriter, r *http.Request) {
	var wt Weather
	var f Forcast
	text := r.URL.Query().Get("text")
	resp, err := http.Get("https://api.weather.gov/points/" + text)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	err = json.Unmarshal([]byte(body), &wt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	resp, err = http.Get(wt.Properties.Forecast)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	err = json.Unmarshal([]byte(body), &f)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	rs := MatterMost.Response{
		ResponseType: "in_channel",
		Text:         "",
	}

	if len(f.Properties.Periods) <= 0 {
		rs.Text = "Unable to load forcast\n"
	} else {
		formattedResult := "| When | Forecast |"
		formattedResult += "\n| :------ | :-------|"
		for i := 0; i < len(f.Properties.Periods); i++ {
			p := f.Properties.Periods[i]
			formattedResult += "\n| " + p.Name + " | " + p.DetailedForecast + " |"
		}

		rs.Text = formattedResult
	}

	toMM, err := json.Marshal(rs)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte(toMM))
}

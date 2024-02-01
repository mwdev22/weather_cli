package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// structurizing response data
type Weather struct {
	// all structures are based on response structure to perfectly fit in and cover important data
	Location struct {
		Name    string `json:"name"`
		Country string `json:"country"`
	} `json:"location"`

	Current struct {
		Temp_C    float64 `json:"temp_c"`
		Condition struct {
			Text string `json:"text"`
		} `json:"condition"`
	} `json:"current"`

	Forecast struct {
		Forecastday []struct {
			Hour []struct {
				Time_Epoch int64   `json:"time_epoch"`
				Temp_C     float64 `json:"temp_c"`
				Condition  struct {
					Text string `json:"text"`
				} `json:"condition"`
				Chance_of_rain float64 `json:"chance_of_rain"`
			} `json:"hour"`
		} `json:"forecastday"`
	} `json:"forecast"`
}

// Funkcja główna programu
func main() {
	// getting data from api
	res, err := http.Get("http://api.weatherapi.com/v1/forecast.json?key=3c633229489f4b71ba8215454242401&q=warsaw&aqi=no")
	if err != nil {
		panic(err)
	}
	defer res.Body.Close() // zwalnianie zasobów

	if res.StatusCode != 200 {
		panic("weather api not available")
	}

	// reading response data to slice of bytes
	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	var weather Weather
	// parsing response data to weather struct
	err = json.Unmarshal(body, &weather)
	if err != nil {
		panic(err)
	}

	location, current, hours := weather.Location, weather.Current, weather.Forecast.Forecastday[0].Hour
	fmt.Printf("Location: %s, Weather: %s: %.0f C degrees, %s\n",
		location.Name,
		location.Country,
		current.Temp_C,
		current.Condition.Text,
	)
	for _, hour := range hours {
		date := time.Unix(hour.Time_Epoch, 0)

		if date.Before(time.Now()) {
			continue
		}

		fmt.Printf("%s - %.0f C degrees, %.0f%% rain possibility, %s\n",
			date.Format("15:04"),
			hour.Temp_C,
			hour.Chance_of_rain,
			hour.Condition.Text,
		)
	}
}

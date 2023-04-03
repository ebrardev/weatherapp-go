package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)

type apiConfigData struct {
	OpenWeatherMapApiKey string `json:"openWeatherMapApiKey"`
}

type weatherData struct {
	Name string `json:"name"`
	Main struct {
		Kelvin float64 `json:"temp"`
	} `json:"main"`
}

func loadApiConfig(filename string) (apiConfigData, error) {

	bytes, err := ioutil.ReadFile(filename)

	if err != nil {
		return apiConfigData{}, err
	}
	var c apiConfigData

	err = json.Unmarshal(bytes, &c)
	if err != nil {
		return apiConfigData{}, err
	}
	return c, nil
}

func hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello! from go\n"))
}
func query(city string) (string, error) {
	apiConfig, err := loadApiConfig(".apiConfig")
	if err != nil {
		return "", err
	}
	resp, err := http.Get("http://api.openweathermap.org/data/2.5/weather?APPID=" + apiConfig.OpenWeatherMapApiKey + "&q=" + city)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var data weatherData
	err = json.Unmarshal(body, &data)
	if err != nil {
		return "", err
	}
	data.Main.Kelvin -= 273.15 // Kelvin to Celsius
	json, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return string(json), nil
}

func main() {

	http.HandleFunc("/hello", hello)
	http.HandleFunc("/weather/",
		func(w http.ResponseWriter, r *http.Request) {
			city := strings.SplitN(r.URL.Path, "/", 3)[2]
			data, err := query(city)

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return

			}
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			json.NewEncoder(w).Encode(data)

		})

	http.ListenAndServe(":8080", nil)

}
